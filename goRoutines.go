package main

/*
void increment() {
}
*/

import "C"

import (
    "fmt"
    "os"
    "net"
    "time"
    "flag"
    )


func errHandler(err error) {

    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}


func chanOut( conn net.Conn, d chan string) {
    var data string
    i := 0

//    fmt.Println("chanOut started")

    for {
        fmt.Println("chanOut ",i)

        data = <- d
        conn.Write ( []byte(data) )
        i++
    }
}

func chanIn( conn net.Conn, d chan byte) {
    buf := make([]byte, 8 )

    i := 0
    idx := 0

    for {
        n, err := conn.Read( buf )
        errHandler(err)

        fmt.Printf("%d:chanIn: ", n )
        fmt.Println( buf )

        d <- byte(n)
        if buf[0] == byte(33) {
            fmt.Println( "POWER FAILURE" )
        }

        for idx=0;idx<n;idx++ {
            d <- buf[idx]
        }
        i++
    }
}

func timer( c chan bool) {
    time.Sleep( 5 * time.Second)
    c <- true
}

func main() {
    //
    // Declare channels
    //
//    C.increment();
    runFlag := true


    toUps   := make(chan string,1)
    fromUps := make(chan byte,8)
    timeout := make(chan bool,1)

    addressPtr := flag.String("address", "192.168.0.143", "a string")
    portPtr := flag.Int("port",4001,"an int")
    verbosePtr := flag.Bool("verbose",false,"a boolean")

    flag.Parse()

    verbose := *verbosePtr
    address := *addressPtr
    port    := *portPtr

    host := fmt.Sprintf("%s:%d",address,port);
    conn, err := net.Dial("tcp", host)
    errHandler( err )

    if verbose {
        fmt.Println("Host    : ",host)
    }

//    time.Sleep( 5 * time.Second)

    fmt.Println("Start go routines")
    go timer(timeout)
    go chanOut( conn, toUps )
    go chanIn( conn, fromUps )

    for runFlag {
        
        toUps <- "Test\n"

        select {
            case <- timeout:
                fmt.Println("Timer expired")
            case <- fromUps:
                d := <- fromUps
                fmt.Println( "data rx:", d )
        }
        time.Sleep( 5 * time.Second)
    }
}
