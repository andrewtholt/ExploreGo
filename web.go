package main

import (
        "fmt"
        "net/http"
        "os"
        "strconv"
        "flag"
        //        "reflect"
       )

type Settings struct {
    debug bool
    xml bool
    yaml bool
    json bool
}

type Ups struct {
    lineVoltage float64
}

func ( u Ups) SetLineVolts(value string) {
    v,err := strconv.ParseFloat(value,32)
    errHndlr(err)

    u.lineVoltage = v
}

func ( u Ups ) Dump() {
    fmt.Printf( "Line Voltage:%4.2f\n", u.lineVoltage )
}

func errHndlr(err error) {
    if err != nil {
        fmt.Println("error:", err)
            os.Exit(1)
    }   
}

func ( u Ups ) FormatOutput() string {
    var out string

    if s.xml {
        out = "This will be XML"
    }

    if s.yaml {
        out = "This will be YAML"
    }

    if s.json {
        out = "This will be JSON"
    }
    return out
}

// Globals

var u Ups
var s Settings

func handler(w http.ResponseWriter, r *http.Request) {

        var mimeType string

        if s.xml {
            mimeType = "text/ups+xml,charset=utf-8"
        }

        if s.yaml {
            mimeType = "text/ups+yaml,charset=utf-8"
        }

        if s.json {
            mimeType = "text/ups+json,charset=utf-8"
        }

        w.Header().Set("Content-Type",mimeType)

        switch r.Method {

            case "GET":
                fmt.Fprintf(w,"%s\n\n", u.FormatOutput())
                
            case "POST", "PUT": 
                fmt.Fprintf(w,"\n\n%s\n\n", r.PostForm)
                fmt.Fprintf(w,"\n\n%s\n\n", r.Form)

                r.ParseForm()

                for key, value := range r.Form {
                    fmt.Fprintf(w,"Key:%s Value:%s\n", key, value)

                        if key == "LINE_VOLTAGE" {
                            u.SetLineVolts( value[0] )
                        }
                }
        }
}

func testHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, Just Testing X%sX", r.URL.Path[1:])
}

func main() {

    debugPtr := flag.Bool("debug",false, "a bool")
    portPtr  := flag.Int("port",8080,"an int")
    xmlPtr   := flag.Bool("xml",true, "a bool")
    yamlPtr   := flag.Bool("yaml",false, "a bool")
    jsonPtr   := flag.Bool("json",false, "a bool")

    flag.Parse()

    if *yamlPtr {
        *xmlPtr=false
        *jsonPtr=false
    }

    if *jsonPtr {
        *xmlPtr=false
        *yamlPtr=false
    }

    if *debugPtr {
        fmt.Println("debug :", *debugPtr)
        fmt.Println("xml   :", *xmlPtr)
        fmt.Println("yaml  :", *yamlPtr)
        fmt.Println("json  :", *jsonPtr)
        fmt.Println("Port  :", *portPtr)
    }

    s.xml  = *xmlPtr
    s.yaml = *yamlPtr
    s.json = *jsonPtr

    http.HandleFunc("/", handler)
    http.HandleFunc("/test", testHandler)

    port := fmt.Sprintf(":%d", *portPtr)

    http.ListenAndServe(port, nil)
}

