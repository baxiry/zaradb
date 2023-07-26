package main

import (
	"dblite"
	"net/http"

	"github.com/go-chi/chi"
)

var engine = dblite.NewEngine()

func main() {

	engine.Run()

	// TODO close programe greatfully
	defer engine.Stop()

	// start network
	router := chi.NewRouter()

	router.Get("/ws", dblite.DemonNet)

	http.ListenAndServe(":1111", router)
}
