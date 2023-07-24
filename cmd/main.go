package main

import (
	"dblite"
	"net/http"

	"github.com/go-chi/chi"
)

var pages dblite.Pages

func main() {
	// run store enginge

	pages = *dblite.NewPages()

	defer pages.Close()

	// start network
	router := chi.NewRouter()

	router.Get("/ws", dblite.Demon)

	http.ListenAndServe(":1111", router)
}
