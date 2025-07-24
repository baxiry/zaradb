package web

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"zaradb/engine"

	"github.com/gorilla/websocket"
)

//go:embed static
var Content embed.FS

type Notify struct {
	msg     string
	msgType int
	err     bool
}

var (
	channel = make(chan Notify, 100)
	note    = Notify{}
)

// render static shell.html file
func Queries(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	query, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result := engine.HandleQueries(string(query))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

// redirect to shell page temporary
func Dev(w http.ResponseWriter, r *http.Request) {
	f, err := Content.ReadFile("static/dev.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(f))
}

// render static shell.html file
func Shell(w http.ResponseWriter, r *http.Request) {
	f, err := Content.ReadFile("static/shell.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(f))
}

// redirect to shell page temporary
func Index(w http.ResponseWriter, r *http.Request) {
	// TODO create index page
	http.Redirect(w, r, "http://localhost:1111/shell", http.StatusSeeOther)
}

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
	var msg []byte

	for {
		if note.err {
			log.Println(err)
			return
		}
		note.msgType, msg, err = c.ReadMessage()
		if err != nil {
			log.Println("Request ERROR! ReadMessage: ", err)
			note.err = true
			channel <- note
			return
		}

		// Hande all of Queries
		note.msg = engine.HandleQueries(string(msg))
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
		err = c.WriteMessage(note.msgType, []byte(note.msg))
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
		msgType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Ws ERROR! ReadMessage: ", err)
			break
		}

		// Hande all of Queries
		result := engine.HandleQueries(string(msg))

		// send result to client
		err = c.WriteMessage(msgType, []byte(result))
		if err != nil {
			log.Println("Ws ERROR! WriteMessage: ", err)
			break
		}
	}

}
