package main

import (
	"dblite"
	"log"

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

	dblite.Server.Handler = router

	err := dblite.Server.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
	}

	log.Println("database is run ...")

	dblite.Shutdown()
}
