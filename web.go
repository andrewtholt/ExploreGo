package main

import (
        "fmt"
        "net/http"
        "os"
        "strconv"
        "flag"
        //        "reflect"
       )

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

var u Ups

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, sent to root %s! and %s", r.URL.Path[1:],r.Method)

        switch r.Method {

            case "GET":
                
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

    http.HandleFunc("/", handler)

        http.HandleFunc("/test", testHandler)
        http.ListenAndServe(":8080", nil)
}

