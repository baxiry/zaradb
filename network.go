package dblite

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO connect with muliple clients

var upgrader = websocket.Upgrader{} // default options

// DemonNet listens incoming queries form ws & send result
func DemonNet(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}
		log.Printf("recv: %s", message)

		// Hande all of Queries
		result := HandleQueries(string(message))

		// send result to client

		err = c.WriteMessage(mt, []byte(result))
		if err != nil {
			panic(err)
		}
	}
}
