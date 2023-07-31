package dblite

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO connect with muliple clients

type Notify struct {
	message     string
	messageType int
}

var Channel = make(chan Notify, 1)

// DemonNet listens incoming queries form ws & send result
func Resever(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	var note Notify
	var message []byte

	for {
		note.messageType, message, err = c.ReadMessage()
		if err != nil {
			fmt.Println("ERROR! :Panic ReadMessage ", err)
			break
		}
		//note.typeMessage = messageType

		//log.Printf("Recve: %s", message)

		// Hande all of Queries
		note.message = HandleQueries(string(message))

		Channel <- note

	}
}

// DemonNet listens incoming queries form ws & send result
func Sender(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	var note Notify

	for {
		note = <-Channel
		// send result to client
		err = c.WriteMessage(note.messageType, []byte(note.message))
		if err != nil {
			fmt.Println("ERROR! :Panic WriteMessage ", err)
			break
		}

	}
}
