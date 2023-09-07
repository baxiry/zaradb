package main

import (
	"net/http"
	db "zaradb/dblite"

	"github.com/go-chi/chi"
)

var engine = db.NewEngine()

func main() {

	engine.Run()

	// TODO close programe greatfully
	defer engine.Stop()

	// start network
	router := chi.NewRouter()
	//	router := http.NewServeMux()

	router.Get("/ws", db.Ws)

	// endpoints for speed network
	router.Get("/query", db.Resever)
	router.Get("/result", db.Sender)

	http.ListenAndServe(":1111", router)
}
