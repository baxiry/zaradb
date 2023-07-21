package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/lesismal/nbio/nbhttp"
)

func main() {
	router := chi.NewRouter()

	router.Get("/ws", onWebsocket)

	svr := nbhttp.NewServer(nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{"localhost:8080"},
	})

	svr.Handler = router

	err := svr.Start()
	if err != nil {
		log.Fatalf("nbio.Start failed: %v\n", err)
	}

	log.Println("database is run ...")

	// shutdown

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	svr.Shutdown(ctx)
}
