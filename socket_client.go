package main

import ( 
        "fmt"
        "net"
        "bufio"
        "strings"
        "reflect"
//        "strconv"
        "os"
        "github.com/fzzy/radix/redis"
        "time"
       )

func errHndlr(err error) {
    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}

func main() {
      c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
      errHndlr(err)
      defer c.Close()

      fmt.Println( reflect.TypeOf(c) )

      host := "192.168.0.143:4001"
      conn, err := net.Dial("tcp", host)

      if err != nil {
          fmt.Println("Network Error")
          os.Exit(1)
      } 

      fmt.Println( reflect.TypeOf(conn) )
      getLineFrequency(conn,c);
      getLineVoltage(conn,c);
      getOutputVoltage(conn,c);
      getBatteryVoltage(conn,c);
      getCauseOfTransfer(conn,c);
      getRunTime(conn,c);

      fmt.Println(err)
}

func getRunTime(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"j\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHndlr(err)

    tmp := strings.TrimSpace(status);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")

    bufio.NewReader(c).ReadString('\n')

    fmt.Printf("Runtime : %s minutes\n",data);

    red.Cmd("set", "RUNTIME", data,"ex","90")
}

func getBatteryVoltage(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"B\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHndlr(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    fmt.Printf("Battery Voltage: %s\n",data);

    r := red.Cmd("set", "BATTERY_VOLTAGE", data,"ex","90")
    errHndlr(r.Err)
}

func getLineFrequency(c net.Conn, red *redis.Client)  {
    fmt.Fprintf(c,"F\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHndlr(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    red.Cmd("set", "LINE_FREQ", data ,"ex","90")
}

func getLineVoltage(c net.Conn, red *redis.Client)  {
    fmt.Fprintf(c,"L\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHndlr(err)

    data := strings.TrimSpace(status);

    bufio.NewReader(c).ReadString('\n')

    red.Cmd("set", "LINE_VOLTS", data,"ex","90")
}

func getOutputVoltage(c net.Conn, red *redis.Client) {
    fmt.Fprintf(c,"O\n")
    status, err := bufio.NewReader(c).ReadString('\n')
    errHndlr(err)

    data := strings.TrimSpace(status);
    bufio.NewReader(c).ReadString('\n')

    red.Cmd("set", "OUTPUT_VOLTS", data,"ex","90")
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
    errHndlr(err)
    data := strings.TrimSpace(status);
//    bufio.NewReader(c).ReadString('\n')

    r := red.Cmd("set", "COT", data, "ex","90")
    errHndlr(r.Err)

    switch status {
    }

    fmt.Printf("COT: %s\n",status);
    return data
}
