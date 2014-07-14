
package ups

import (
    "os"
    "net"
    "fmt"
    "bufio"
    "strings"
)

type UpsData struct {
    c net.Conn
    batteryVoltage string
    batteryLevel string
    lineVoltage string
    runTime string
}

func (u UpsData) Dump() {
    fmt.Println("\tData Dump")
    fmt.Println("\t==== ====")

    fmt.Println("Connection: ",u.c)
    fmt.Println("Battery   : ",u.batteryVoltage)
    fmt.Println("Line Volts: ",u.lineVoltage)
    fmt.Println("Runtime   : ",u.runTime)
    fmt.Println("\t=========")
}

func (u *UpsData) UpdateBatteryLevel() error {
    fmt.Fprintf(u.c,"f\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')

    fmt.Println("runtime is ", status)
    tmp := strings.TrimSpace(status);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")
    fmt.Println(data)

    bufio.NewReader(u.c).ReadString('\n')

    u.batteryLevel = data

    return err
}

func (u *UpsData) UpdateRuntime() error {
    fmt.Fprintf(u.c,"j\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')

    fmt.Println("runtime is ", status)
    tmp := strings.TrimSpace(status);
    data := strings.TrimLeft(strings.TrimRight(tmp,":"),"0")
    fmt.Println(data)

    bufio.NewReader(u.c).ReadString('\n')

    u.runTime = data

    return err
}

func (u *UpsData) UpdateBatteryVoltage() error {
    fmt.Fprintf(u.c,"B\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    data := strings.TrimSpace(status);
    fmt.Println(data)

    bufio.NewReader(u.c).ReadString('\n')

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
    fmt.Fprintf(u.c,"L\n")
    status, err := bufio.NewReader(u.c).ReadString('\n')
    errHandler(err)

    data := strings.TrimSpace(status);

    bufio.NewReader(u.c).ReadString('\n')

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

    conn, err := net.Dial("tcp", source )
    errHandler( err )

    p.c = conn

    /*
    errHandler(p.UpdateBatteryVoltage())
    errHandler(p.UpdateLineVoltage())
    */

    return p
}

