/*
    Class constructor/initialiser example.
*/
package main

import (
    "fmt"
)

type Thing struct {
    Name string
    Num int
}

func Create(t string) *Thing {
    p := new(Thing)

    p.Name = t
    p.Num = 42

    return p
}

func (u Thing)Dump() {
    fmt.Println("Name is ",u.Name)
}

func (u *Thing) Set(n string) {
    u.Name = n
}


func main() {
    x := Create("Andrew")
    x.Dump()

    x.Set("Fred")
    x.Dump()

}
