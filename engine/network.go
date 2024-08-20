package engine

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Notify struct {
	message     string
	messageType int
	err         bool
}

var (
	channel = make(chan Notify, 100)
	note    = Notify{}
)

// Request listens incoming queries form client
func Request(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ERROR! upgrade connection ", err)
		return
	}
	defer c.Close()

	var note Notify
	var message []byte

	for {
		if note.err {
			log.Println(err)
			return
		}
		note.messageType, message, err = c.ReadMessage()
		if err != nil {
			log.Println("Request ERROR! ReadMessage: ", err)
			note.err = true
			channel <- note
			return
		}

		// Hande all of Queries
		note.message = HandleQueries(string(message))
		channel <- note
	}
}

// Response sends results to clients
func Response(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ERROR! upgrade connection: ", err)

		return
	}
	defer c.Close()

	for {
		note = <-channel
		if note.err {
			log.Println(err)
			return
		}
		// send result to client
		err = c.WriteMessage(note.messageType, []byte(note.message))
		if err != nil {
			log.Println("Response ERROR! WriteMessage: ", err)
			note.err = true
			channel <- note
			return
		}
	}
}

// Ws listens incoming queries form ws & send result
func Ws(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Ws ERROR! ReadMessage: ", err)
			break
		}

		// Hande all of Queries
		result := HandleQueries(string(message))

		// send result to client
		err = c.WriteMessage(messageType, []byte(result))
		if err != nil {
			log.Println("Ws ERROR! WriteMessage: ", err)
			break
		}
	}

}
