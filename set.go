package main

import (
    "fmt"
    //"time"
    "flag"
    //"strconv"

	"github.com/nats-io/go-nats"
)


func main() {
    helps := `
        -h      to show tis help message
        -t      to set time duration by second e.g: -tm 5 this add 5 second bitween every messages
        -s      to set subject/topic
        // another setting here
        // another setting here

    `
    // connect to brocker
    n, err := nats.Connect("0.0.0.0:4222")
    if err != nil {
        fmt.Println("no connect", err)
    }
    defer n.Close()


    //var topic string
	var timeUp string
    var topic string
    var help string
    //timeUp := 3

    flag.StringVar(&topic, "s", "", "specify of topic")
    flag.StringVar(&help, "h", "", "show help message")
    flag.StringVar(&timeUp, "t", "3", "adding a wait time by second bitween messages")

    flag.Parse()

    n.Publish("setting", []byte(timeUp))


    fmt.Println("is update")
}
