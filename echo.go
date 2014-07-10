package main

import (
        "flag"
        "fmt"
        "strings"
       )

func main() {
    var omitNewline bool
        flag.BoolVar(&omitNewline, "n", false, "don't print final newline")
        flag.Parse() // Scans the arg list and sets up flags.

        str := strings.Join(flag.Args(), " ")
        if omitNewline {
            fmt.Print(str)
        } else {
            fmt.Println(str)
        }
}
