package dblite

import (
	"fmt"
	"net/http"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

// Demon listens incoming queries form ws & send result
func Demon(resp http.ResponseWriter, req *http.Request) {
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
}
