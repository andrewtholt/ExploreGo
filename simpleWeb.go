package main

import (
    "fmt"
    "net/http"
    //    "os"
    //    "strconv"
    "flag"
    //    "strings"
    //    "github.com/fzzy/radix/redis"
    //    "time"
    //    "reflect"
)

const XML int  = 1
const YAML int = 2
const JSON int = 3

var outputType int

func testHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, Just Testing X%sX\n\n", r.URL.Path[1:])
}

func upsHandler(w http.ResponseWriter, r *http.Request) {

    switch outputType {
        case XML:
            fmt.Fprintf(w, "<?xml version=\"1.0\"?>\n")
            fmt.Fprintf(w, "<UPS>\n");
            fmt.Fprintf(w, "</UPS>\n");

        case YAML:
        case JSON:
    }
    /*
    fmt.Fprintf(w, "Hi there, Just Handling X%sX\n\n", r.URL.Path[1:])
    fmt.Fprintf(w, "Output Type is : %d\n\n",outputType)

    fmt.Fprintf(w, "Hi there, Just Handling X%sX\n\n", r.URL.Path[1:])
    fmt.Fprintf(w, "Output Type is : %d\n\n",outputType)
    */
}


func main() {

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

    http.HandleFunc("/ups", upsHandler)
    http.HandleFunc("/test", testHandler)

    //    port := fmt.Sprintf(":%d", *portPtr)

    http.ListenAndServe(port, nil)
}
