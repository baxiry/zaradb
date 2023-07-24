package dblite

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options
// Demon listens incoming queries form ws & send result
func Demon(w http.ResponseWriter, r *http.Request) {
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

		result := HandleQueries(string(message))

		// send result to client

		err = c.WriteMessage(mt, []byte(result))
		if err != nil {
			panic(err)
		}
	}
}

/*

	upgrader := websocket.NewUpgrader()
	upgrader.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		result := HandleQueries(string(data))

		// send result to client
		conn.WriteMessage(messageType, []byte(result))

		fmt.Println("data reseved : ", string(data))
	})

	upgrader.OnOpen(func(conn *websocket.Conn) {
		fmt.Println("OnOpen:", conn.RemoteAddr().String())
	})

	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		panic(err)
	}

	conn.OnClose(func(c *websocket.Conn, err error) {
		fmt.Println("OnClose:", c.RemoteAddr().String(), err)
	})
*/
