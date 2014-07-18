
package main

import (
        "fmt"
        "os"
        "net"
        "time"
        "flag"
        "reflect"
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
        fmt.Println("outToUps:", data)

        conn.Write( []byte(data ))
    }

}

func inFromUps ( conn net.Conn, c chan byte) {
    
    localBuffer := make( []byte,1 )
    for {
        n, err := conn.Read( localBuffer )
        fmt.Println("inFromUps : ", n)

        if localBuffer[0] != 10 {
            errHandler( err )
            c <- localBuffer[0]
        }
    }
}

func getLineVoltage(r *redis.Client, t chan string, f chan byte) string {
    var c byte
    buffer := make([]byte,8)
    t <- "L\n"

    for i:=0 ; i< 5 ; i++ {
        c = <- f
        fmt.Println( c )
        buffer[i] = c
    }

    data := strings.Trim(string(buffer),"\x00");
    r.Cmd("set","LINE_VOLTAGE", data,"ex","90")
    return data
}

func main() {
    toUps   := make(chan string,1)
    fromUps := make(chan byte,1)

//    timeout := make(chan bool,1)

    addressPtr := flag.String("address", "192.168.0.143", "a string")
    redisPtr := flag.String("redis", "127.0.0.1", "a string")
    portPtr := flag.Int("port",4001,"an int")
    verbosePtr := flag.Bool("verbose",false,"a boolean")
    delayPtr   := flag.Int("delay",60,"an int")

    flag.Parse()

    verbose      := *verbosePtr
    address      := *addressPtr
    port         := *portPtr
    redisAddress := *redisPtr
    delay        := *delayPtr

    host := fmt.Sprintf("%s:%d", address, port)
    redisHost := fmt.Sprintf("%s:6379",redisAddress)

    if verbose {
        fmt.Println("Serial Server : ", host )
        fmt.Println("Serial Port   : ", host )
        fmt.Println("Redis  Server : ", redisHost )
        fmt.Println("Delay         : ", delay, "Seconds" )
    }

    redisConn, err := redis.DialTimeout("tcp", redisHost, time.Duration(10)*time.Second)
    errHandler( err )

    fmt.Println( reflect.TypeOf(redisConn) )

    conn, err := net.Dial("tcp", host)
    errHandler( err )

    go outToUps(conn, toUps )
    go inFromUps(conn, fromUps )

    for {
        fmt.Println( getLineVoltage(redisConn, toUps, fromUps))

        fmt.Println("END")
        time.Sleep( time.Duration(delay) * time.Second )
    }
}

