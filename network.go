package main

import (
	"log"
	"net/http"
	"zaradb/store"

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
		log.Print("when upgrade ", err)
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
			log.Println("ReadMessage ", err)
			note.err = true
			Channel <- note
			return
		}

		// Hande all of Queries
		note.message = store.HandleQueries(string(message))

		Channel <- note

	}
}

// Sender  sends results to clients
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
		if note.err {
			log.Println(err)
			return
		}
		// send result to client
		err = c.WriteMessage(note.messageType, []byte(note.message))
		if err != nil {
			log.Println("ERROR! :Panic WriteMessage ", err)

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
			log.Println("ERROR! :Panic ReadMessage ", err)

			break
		}

		// Hande all of Queries
		//start := time.Now()
		result := store.HandleQueries(string(message)) // + "\n" + time.Since(start).String()

		// send result to client
		err = c.WriteMessage(messageType, []byte(result))
		if err != nil {
			log.Println("ERROR! :Panic WriteMessage ", err)
			break
		}
	}

}
