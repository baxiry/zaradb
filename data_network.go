package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

func onWebsocket(resp http.ResponseWriter, req *http.Request) {
	upgrader := websocket.NewUpgrader()
	upgrader.OnMessage(func(conn *websocket.Conn, messageType websocket.MessageType, data []byte) {
		// echo
		conn.WriteMessage(messageType, data)
		fmt.Println("data reseved : ", string(data))
	})
	upgrader.OnOpen(func(conn *websocket.Conn) {
		log.Println("OnOpen:", conn.RemoteAddr().String())
	})

	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		panic(err)
	}
	conn.OnClose(func(c *websocket.Conn, err error) {
		log.Println("OnClose:", c.RemoteAddr().String(), err)
	})
}
