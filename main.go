// Zaradb lite faset document database
package main

import (
	"net/http"
	db "zaradb/dblite"
)

var engine = db.NewEngine()

func main() {
	// TODO close programe greatfully.

	engine.Run()
	defer engine.Stop()

	// standard endpoint
	http.HandleFunc("/ws", Ws)

	// endpoints for speed network
	http.HandleFunc("/query", Resever)
	http.HandleFunc("/result", Sender)

	http.ListenAndServe(":1111", nil)
}
