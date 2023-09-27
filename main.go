package main

import (
	"net/http"
	db "zaradb/dblite"
)

var engine = db.NewEngine()

func main() {
	// TODO close programe greatfully

	engine.Run()
	defer engine.Stop()

	// standard endpoint
	http.HandleFunc("/ws", db.Ws)

	// endpoints for speed network
	http.HandleFunc("/query", db.Resever)
	http.HandleFunc("/result", db.Sender)

	http.ListenAndServe(":1111", nil)
}
