package main

import ( 
        "fmt"
        "net"
        "bufio"
        "strings"
//        "reflect"
        "os"
        "github.com/fzzy/radix/redis"
        "time"
        "flag"
       )

func errHandler(err error) {
    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}

const OK=0
const POWER_FAIL=1
const POWER_RETURN=2


func myRead(c net.Conn, len int) (string, int) {
    buf := make([]byte, 16)
    out := make([]byte, 16)
    i := 0
    fail :=0
    lfCount := 0

    limit := len-1
    n,err := c.Read( buf )
    errHandler(err)

    switch buf[0] {
        case '!':
            return "",POWER_FAIL
        case '$':
            return "",POWER_RETURN
            
    }

    out[i] = buf[0]
    i=i+1


    fmt.Println("Limit=",limit)
    for i < limit {
        n,err = c.Read( buf )
        errHandler(err)

        out[i]=buf[0]
        if 10 == buf[0] {
            lfCount++
        }

        if lfCount >= 2 {
            out[i] = 0x00
            break;
        }

        if 0 == lfCount {
            i=i+1
        }

        fmt.Println("Read ",n)
        fmt.Println("i=",i)
        fmt.Println( buf )
        fmt.Println( out )

        if err != nil || i > limit { break }
    }

    fmt.Println("Broke")
    return string(out[:len]), fail
}

func main() {
    addressPtr := flag.String("address", "192.168.0.143", "a string")
    portPtr := flag.Int("port",4001,"an int")

    flag.Parse()


    /*
        Connect to redis
    */
    redisHost := "127.0.0.1:6379"
    c, err := redis.DialTimeout("tcp", redisHost, time.Duration(10)*time.Second)
    errHandler(err)
    defer c.Close()

    /*
        Connect to serialServer.
    */
    host := fmt.Sprintf("%s:%d",*addressPtr,*portPtr)
    fmt.Println("Connect to ",host)

    conn, err := net.Dial("tcp", host)

    if err != nil {
        fmt.Println("33:Network Error")
        os.Exit(1)
    } 

    fmt.Println("Here ... ")
    getLineVoltage(conn,c);
//    getLineFrequency(conn,c);
//    getOutputVoltage(conn,c);
//    getBatteryVoltage(conn,c);
//    getBatteryLevel(conn,c);
//    getCauseOfTransfer(conn,c);
//    getRunTime(conn,c);

    fmt.Println("... and here.")
    fmt.Println(err)
}

func getRunTime(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"j\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    tmp := strings.TrimSpace(status);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")

    bufio.NewReader(c).ReadString('\n')

    fmt.Printf("Runtime : %s minutes\n",data);

    red.Cmd("set", "RUNTIME", data,"ex","90")
}

func getBatteryLevel(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"f\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    fmt.Printf("Battery Level  : %s\n",data);

    r := red.Cmd("set", "BATTERY_LEVEL", data,"ex","90")
    errHandler(r.Err)
}

func getBatteryVoltage(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"B\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    fmt.Printf("Battery Voltage: %s\n",data);

    r := red.Cmd("set", "BATTERY_VOLTAGE", data,"ex","90")
    errHandler(r.Err)
}

func getLineFrequency(c net.Conn, red *redis.Client)  {
    fmt.Fprintf(c,"F\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    red.Cmd("set", "LINE_HZ", data ,"ex","90")
}

func getLineVoltage(c net.Conn, red *redis.Client)  {
    fmt.Fprintf(c,"L\n")

    status,fail := myRead(c,5) 
    fmt.Println("Status ",fail)

    if  fail == OK {
        data := strings.TrimSpace(status);
        fmt.Println("Data ",data)

    /*
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);

    bufio.NewReader(c).ReadString('\n')
    */

//       r := red.Cmd("set", "LINE_VOLTAGE", data,"ex","90")
       r := red.Cmd("set", "LINE_VOLTAGE", data)
       if r.Err != nil {
        fmt.Println("Redis error.")
        errHandler(r.Err)
       }

       r = red.Cmd("expire", "LINE_VOLTAGE", 90)
       if r.Err != nil {
        fmt.Println("Redis error.")
        errHandler(r.Err)
       }
    }
}

func getOutputVoltage(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"O\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    red.Cmd("set", "OUTPUT_VOLTAGE", data,"ex","90")
}

/*
    Cause Of Transfer Codes.

    R = unacceptable utility voltage rate of change, 
    H = high utility voltage, 
    L = low utility voltage, 
    T = line voltage notch or spike, 
    O = no transfers yet (since turnon), 
    S = transfer due to serial port U command or activation of UPS test from front panel, 
    NA = transfer reason still not available (read again).
*/

func getCauseOfTransfer(c net.Conn, red *redis.Client) string {
    fmt.Println("Cause");
    fmt.Fprintf(c,"G\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHandler(err)
    data := strings.TrimSpace(status);

    r := red.Cmd("set", "COT", data, "ex","90")
    errHandler(r.Err)

    switch status {
    }

    fmt.Printf("COT: %s\n",status);
    return data
}
