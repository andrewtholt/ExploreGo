package main

import ( 
        "fmt"
//        "net"
//        "bufio"
//        "strings"
//        "reflect"
//        "strconv"
        "os"
//        "container/list"
        "github.com/fzzy/radix/redis"
        "time"
       )

func errHndlr(err error) {
    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}

func head() {
    fmt.Println("<?xml version=\"1.0\"?>")
    fmt.Println("<UPS>")
}

func tail() {
    fmt.Println("</UPS>")
}

func main() {
    c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
    errHndlr(err)
    defer c.Close()

    keys := []string{"BATTERY_VOLTS","LINE_VOLTS","OUTPUT_VOLTS","COT"}

    ls,err := c.Cmd("mget",keys).List()

    head()
    for index, element := range ls {
        fmt.Println("\t<" + keys[index] + ">" + element + "</" + keys[index] + ">")
//        fmt.Println("\t<value>" + element + "</value>")
    }

    r := c.Cmd("get", "BATTERY_VOLTS")
    errHndlr(r.Err)

    tail()
}
