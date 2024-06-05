package engine

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO connect with muliple clients

type Notify struct {
	message     string
	messageType int
	err         bool
}

var Channel = make(chan Notify, 100)

// Resever listens incoming queries form clients
func Resever(w http.ResponseWriter, r *http.Request) {

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
			log.Println("ERROR! ReadMessage: ", err)
			note.err = true
			Channel <- note
			return
		}

		// Hande all of Queries
		note.message = HandleQueries(string(message))

		Channel <- note

	}
}

// Sender  sends results to clients
func Sender(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{} // default options

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ERROR! upgrade connection: ", err)

		return
	}
	defer c.Close()

	var note Notify

	for {
		note = <-Channel
		if note.err {
			log.Println(err)
			return
		}
		// send result to client
		err = c.WriteMessage(note.messageType, []byte(note.message))
		if err != nil {
			log.Println("ERROR! WriteMessage: ", err)

			note.err = true
			Channel <- note
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
			log.Println("ERROR! ReadMessage: ", err)
			print("ok")

			break
		}

		// Hande all of Queries
		result := HandleQueries(string(message))

		// send result to client
		err = c.WriteMessage(messageType, []byte(result))
		if err != nil {
			log.Println("ERROR! WriteMessage: ", err)
			break
		}
	}

}
