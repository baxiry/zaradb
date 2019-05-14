package main

import (
    "fmt"
    "time"
    "strconv"

	"github.com/nats-io/go-nats"
)


func main() {
    // connect to brocker
    n, err := nats.Connect("0.0.0.0:4222")
    if err != nil {
        fmt.Println("no connect", err)
    }
    defer n.Close()

    fmt.Println("is connected")

    tm := 3
    n.Subscribe("setting", func(m *nats.Msg) {
        tm , _ = strconv.Atoi(string(m.Data))
        fmt.Printf("Received a message: %s\n", string(m.Data))
        fmt.Println("Setting of is changed to ", string(m.Data), "second")
    })


    i := 0
    for {
        tm = tm

        time.Sleep(time.Second*time.Duration(tm))

        n.Publish("room", []byte("hello world "+strconv.Itoa(i))); i++
    }

    time.Sleep(time.Second)



    //n.Publish("a room", []byte("shot down"))


    fmt.Println("Done")
}

