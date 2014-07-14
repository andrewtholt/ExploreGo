
package ups

import (
    "os"
//    "net"
    "fmt"
//    "bufio"
    "strings"
    "github.com/fzzy/radix/redis"
    "time"
)

type UpsData struct {
//    c net.Conn
    batteryVoltage string
    batteryLevel string
    lineVoltage string
    runTime string
    redisConnect *redis.Client
}

func (u UpsData) Dump() {
    fmt.Println("\tData Dump")
    fmt.Println("\t==== ====")

//    fmt.Println("Connection: ",u.c)
    fmt.Println("Battery   : ",u.batteryVoltage)
    fmt.Println("Line Volts: ",u.lineVoltage)
    fmt.Println("Runtime   : ",u.runTime)
    fmt.Println("\t=========")
}

func (u *UpsData) UpdateBatteryLevel() error {

    r,err := u.redisConnect.Cmd("GET","BATTERY_LEVEL").Str()
    fmt.Println("BATTERY LEVEL is ", r)

    u.batteryLevel = r

    /*
    fmt.Fprintf(u.c,"f\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    fmt.Println("runtime is ", status)
    tmp := strings.TrimSpace(status);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")
    fmt.Println(data)

    bufio.NewReader(u.c).ReadString('\n')

    u.batteryLevel = data
    */

    return err
}

func (u *UpsData) UpdateRuntime() error {
    r,err := u.redisConnect.Cmd("GET","RUNTIME").Str()
    fmt.Println("RUNTIME  is ", r)

    tmp := strings.TrimSpace(r);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")

    /*
    fmt.Fprintf(u.c,"j\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    fmt.Println(data)
    bufio.NewReader(u.c).ReadString('\n')
    */

    u.runTime = data

    return err
}

func (u *UpsData) UpdateBatteryVoltage() error {
    /*
    fmt.Fprintf(u.c,"B\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    */

    r,err := u.redisConnect.Cmd("GET","BATTERY_VOLTAGE").Str()
    fmt.Println("BATTERY_VOLTAGE  is ", r)
    data := strings.TrimSpace(r);
    fmt.Println(data)

//    bufio.NewReader(u.c).ReadString('\n')

    u.batteryVoltage = data

    return err
}

func (u *UpsData) GetBatteryLevel() string {
    return u.batteryLevel
}

func (u *UpsData) GetRuntime() string {
    return u.runTime
}

func (u *UpsData) GetBatteryVoltage() string {
    return u.batteryVoltage
}

func (u *UpsData) GetLineVoltage() string {
    return u.lineVoltage
}

func (u *UpsData) UpdateLineVoltage() error {
    r,err := u.redisConnect.Cmd("GET","LINE_VOLTAGE").Str()
    fmt.Println("LINE VOLTAGE  is ", r)

    /*
    fmt.Fprintf(u.c,"L\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    errHandler(err)
    */

    data := strings.TrimSpace(r);

//    bufio.NewReader(u.c).ReadString('\n')

    u.lineVoltage = data

    return err
}

func errHandler(err error) {
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }   
}

func Create( source string ) *UpsData {
    fmt.Printf("Create UPS instance :%s\n",source)
    p := new(UpsData)

    /*
    conn, err := net.Dial("tcp", source )
    errHandler( err )

    p.c = conn
    */

    p.redisConnect, _ = redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
//    errHandler( err )


    /*
    errHandler(p.UpdateBatteryVoltage())
    errHandler(p.UpdateLineVoltage())
    */

    return p
}

