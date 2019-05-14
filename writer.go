
package main

import (
    "fmt"
    "time"
    "log"
    "gopkg.in/mgo.v2"
	"github.com/nats-io/go-nats"
    // "gopkg.in/mgo.v2/bson"
)

type Message struct {
        // Name string
        Data string
}

func main() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
           panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)


    // connect to brocker
    n, err := nats.Connect("0.0.0.0:4222")
    if err != nil {
        fmt.Println("no connect", err)
    } else {
        fmt.Println("is connected")
    }
    defer n.Close()




    n.Subscribe("room", func(m *nats.Msg){

        c := session.DB("nats").C("test")
        err = c.Insert(&Message{string(m.Data)})
        if err != nil {
                log.Fatal(err)
        }

        fmt.Println(string(m.Data))
    })

    for {
        time.Sleep(time.Minute)
    }

}


