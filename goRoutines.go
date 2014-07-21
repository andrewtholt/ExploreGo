
package main

import (
    "fmt"
    "os"
    "net"
    "time"
    "flag"
    //        "reflect"
    "github.com/fzzy/radix/redis"
    "strings"
)


func errHandler(err error) {
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }   
}

func Out() {
}

func In() {
}

func outToUps ( conn net.Conn, c chan string) {

    var data string

    for {
        data = <-c
        conn.Write( []byte(data ))
    }

}

func inFromUps ( conn net.Conn, c chan byte) {
    
    localBuffer := make( []byte,1 )
    for {
        _, err := conn.Read( localBuffer )

        if localBuffer[0] != 10 {
            errHandler( err )
            c <- localBuffer[0]
        }
    }
}

func getCauseOfTransfer(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "G\n"
    c = <- f
    buffer[0] = c

    data := ""

    switch c {
        case 'R':
            data = "ROC"
        case 'H':
            data = "HIGH"
        case 'L':
            data = "LOW"
        case 'T':
            data = "SPIKE"
        case 'O':
            data = "NONE"
        case 'S':
            data = "USER"
        case 'N':
            c = <- f
            data = "NA"
        default:
            data = "ERROR"
    }

    r.Cmd("set","COT", data,"ex","90")
    return data
}

func getLineFrequency(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "F\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","LINE_HZ", data,"ex","90")
    return data
}

func getBatteryVoltage(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "B\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","BATTERY_VOLTAGE", data,"ex","90")
    return data
}

func getBatteryLevel(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "f\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","BATTERY_LEVEL", data,"ex","90")
    return data
}

func getOutputVoltage(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "O\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","OUTPUT_VOLTAGE", data,"ex","90")
    return data
}

func getLineVoltage(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "L\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","LINE_VOLTAGE", data,"ex","90")
    return data
}

func usage() {
    fmt.Println("Usage: goRoutines -help|-verbose|-address <service address>|-port <port number>")
    fmt.Println("\t-redis <redis server address>| -delay <n>\n")

    fmt.Println("If run with no switches the default is as if you entered:")
    fmt.Println("goRutines.go -address 192.168.0.143 -port 4001 -redis 127.0.0.1 -delay 60")

}


func main() {
    toUps   := make(chan string,1)
    fromUps := make(chan byte,1)

    helpPtr    := flag.Bool("help",false,"a boolean")
    addressPtr := flag.String("address", "192.168.0.143", "a string")
    redisPtr   := flag.String("redis", "127.0.0.1", "a string")
    portPtr    := flag.Int("port",4001,"an int")
    verbosePtr := flag.Bool("verbose",false,"a boolean")
    delayPtr   := flag.Int("delay",60,"an int")

    flag.Parse()

    verbose      := *verbosePtr
    address      := *addressPtr
    port         := *portPtr
    redisAddress := *redisPtr
    delay        := *delayPtr
    help         := *helpPtr

    host := fmt.Sprintf("%s:%d", address, port)
    redisHost := fmt.Sprintf("%s:6379",redisAddress)

    if help {
        usage()
        os.Exit(0)
    }
    if verbose {
        fmt.Println("Serial Server : ", address )
        fmt.Println("Serial Port   : ", port )
        fmt.Println("Serial Host   : ", host )
        fmt.Println("================ ")
        fmt.Println("Redis  Server : ", redisHost )
        fmt.Println("Delay         : ", delay, "Seconds" )
    }

    redisConn, err := redis.DialTimeout("tcp", redisHost, time.Duration(10)*time.Second)
    errHandler( err )

    conn, err := net.Dial("tcp", host)
    errHandler( err )

    go outToUps(conn, toUps )
    go inFromUps(conn, fromUps )

    var lineVoltage string
    var lineFrequency string
    var outputVoltage string
    var batteryLevel string
    var batteryVoltage string
    var causeOfTransfer string

    for {
        lineVoltage = getLineVoltage(redisConn, toUps, fromUps)
        lineFrequency = getLineFrequency(redisConn, toUps, fromUps)
        outputVoltage =  getOutputVoltage(redisConn, toUps, fromUps)
        batteryLevel = getBatteryLevel(redisConn, toUps, fromUps)
        batteryVoltage = getBatteryVoltage(redisConn, toUps, fromUps)
        causeOfTransfer = getCauseOfTransfer(redisConn, toUps, fromUps)

        if verbose {
            fmt.Println( "Line Voltage     :",lineVoltage,"Volts" )
            fmt.Println( "Line Frequency   :",lineFrequency,"Hz" )
            fmt.Println( "Output Voltage   :",outputVoltage,"Volts" )
            fmt.Println( "Battery Level    :",batteryLevel,"%" )
            fmt.Println( "Battery Voltage  :",batteryVoltage,"Volts" )
            fmt.Println( "Cause Of Transfer:", causeOfTransfer)

            fmt.Println("END")
            time.Sleep( time.Duration(delay) * time.Second )
        }
    }
}

