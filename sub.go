
package main

import (
    "fmt"
    "time"

	"github.com/nats-io/go-nats"
)

func main() {
    // connect to brocker
    n, err := nats.Connect("0.0.0.0:4222")
    for i := 0; i < 500; i++ {
        n , _ := nats.Connect("0.0.0.0:4222")
        defer n.Close()
    }
    if err != nil {
        fmt.Println("no connect", err)
    }
    defer n.Close()
    fmt.Println("is connected")


    n.Subscribe("room", func(m *nats.Msg){
        fmt.Println("Receved a message: ", string(m.Data))

            //if string(m.Data) == "shot down" {fmt.Println("publisher is finish")}
    })
    for {
       time.Sleep(time.Second*10)
    }
}



