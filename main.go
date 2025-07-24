// Zaradb lite fast document database
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	eng "zaradb/engine"
	"zaradb/web"
)

var server = &http.Server{
	Addr: Port,
}

func main() {

	db := eng.NewDB("test.db")
	if db == nil {
		log.Fatal("no db")
		return
	}
	defer db.Close()

	fmt.Printf("interacte with zaradb through %s%s\n", Host, Port)

	http.Handle("/static/", http.FileServer(http.FS(web.Content)))

	http.HandleFunc("/", web.Index)
	http.HandleFunc("/index", web.Shell)

	// endpoints for low loading
	http.HandleFunc("/queries", web.Queries)

	http.HandleFunc("/shell", web.Shell)

	// endpoints for high loading
	http.HandleFunc("/query", web.Request)
	http.HandleFunc("/result", web.Response)

	// for pages under development
	http.HandleFunc("/dev", web.Dev)

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down zaradb...")

	// Create a context with timeout to allow active requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("zaradb forced to shutdown:", err)
	}

	log.Println("zaradb exiting")

}
