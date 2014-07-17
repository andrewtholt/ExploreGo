
package main

import (
    "fmt"
    "os"
    "net"
    "time"
    )

var event chan byte
var toUps chan string
var fromUps chan byte
var runFlag bool

func errHandler(err error) {
    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}

func boss(conn net.Conn, c chan byte ) {
}

func iFaceRead(conn net.Conn, d chan byte) {
    run := true

    fmt.Println("iFaceRead")

    buf := make([]byte, 16)
    out := make([]byte, 8)

    for {
    for i :=0;run; {
        n,err := conn.Read( buf)
        errHandler(err)

        fmt.Println("Data length is ",n)

        fmt.Println("Data        is ",buf)
        fmt.Println("Index       is ",i)
        out[i] = buf[0]
        fmt.Println("Out         is ",out)
        i++;

        if buf[0] == 10 { run = false }
    }
    fmt.Println("iFaceRead done")

    d <- buf[0]
    }

}

func iFaceWrite(conn net.Conn, d chan string) {
    fmt.Println("iFaceWrite")

    var buf string

    buf = <- d

    fmt.Println("iFaceOut Data   is ",buf)
    fmt.Println("iFaceOut length is ",len(buf))
    conn.Write([]byte(buf))
    
}

func main() {
    event   = make(chan byte,1)
    toUps   = make(chan string,1)
    fromUps = make(chan byte,1)

    runFlag = true

    host := "192.168.0.143:4001"
    conn, err := net.Dial("tcp", host)

    errHandler(err)

    go iFaceWrite( conn, toUps )
//    conn.Write([]byte("L\r"))
    time.Sleep(1 * time.Second)
    toUps <- "L\r"

    go iFaceRead( conn, fromUps )

    time.Sleep(1 * time.Second)
}

