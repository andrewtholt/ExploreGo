package main

import (
    "fmt"
    "net/http"
    "os"
    //    "strconv"
    "flag"
    //    "strings"
//    "github.com/fzzy/radix/redis"
//    "time"
    "apc_ups"
    "reflect"
)

const XML int  = 1
const YAML int = 2
const JSON int = 3

var upsInstance *ups.UpsData

var outputType int

func errHandler(err error) {
    if err != nil {
        fmt.Println("error:", err)
        os.Exit(1)
    }
}
func testHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, Just Testing X%sX\n\n", r.URL.Path[1:])
}

func upsHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println("x is ",upsInstance)
    upsInstance.UpdateBatteryVoltage()
    upsInstance.UpdateLineVoltage()
    upsInstance.UpdateRuntime()
    upsInstance.UpdateBatteryLevel()

    upsInstance.Dump()


//    c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10)*time.Second)
//    errHandler(err)
//    defer c.Close()

    switch outputType {
        case XML:
            fmt.Fprintf(w, "<?xml version=\"1.0\"?>\n")
            fmt.Fprintf(w, "<UPS>\n");
            fmt.Fprintf(w, "    <BATTERY_VOLTAGE>%s</BATTERY_VOLTAGE>\n", upsInstance.GetBatteryVoltage())
            fmt.Fprintf(w, "    <BATTERY_LEVEL>%s</BATTERY_LEVEL>\n", upsInstance.GetBatteryLevel())

            fmt.Fprintf(w, "    <LINE_VOLTAGE>%s</LINE_VOLTAGE>\n", upsInstance.GetLineVoltage())
            fmt.Fprintf(w, "    <RUNTIME>%s</RUNTIME>\n", upsInstance.GetRuntime())
            fmt.Fprintf(w, "</UPS>\n");

        case YAML:
        case JSON:
    }
}


func main() {
    upsInstance = nil

    upsInstance = ups.Create( "192.168.0.143:4001" )
    fmt.Println( reflect.TypeOf(upsInstance) )

    outputType = XML

    debugPtr := flag.Bool("debug",false, "a bool")
    portPtr  := flag.Int("port",8080,"an int")
    xmlPtr   := flag.Bool("xml",true, "a bool")
    yamlPtr  := flag.Bool("yaml",false, "a bool")
    jsonPtr  := flag.Bool("json",false, "a bool")

    flag.Parse()

    if *xmlPtr {
        outputType = XML
    }

    if *yamlPtr {
        outputType = YAML
    }

    if *jsonPtr {
        outputType = JSON
    }

    port := fmt.Sprintf(":%d",*portPtr)

    if *debugPtr {
        fmt.Println("Debug :", *debugPtr )
        fmt.Println("Port  :", *portPtr )
    }


    //    port := fmt.Sprintf(":%d", *portPtr)

    http.HandleFunc("/ups", upsHandler)
    http.HandleFunc("/test", testHandler)
    http.ListenAndServe(port, nil)
}
