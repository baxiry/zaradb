// Zaradb lite fast document database
package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"zaradb/engine"
)

//go:embed static
var content embed.FS

// TODO: Close program gracefully.
func main() {

	db := engine.NewDB("test.db")
	if db == nil {
		log.Fatal("no db")
		return
	}
	defer db.Close()

	fmt.Printf("interacte with zaradb through %s:%s\n", Host, Port)

	http.Handle("/static/", http.FileServer(http.FS(content)))

	http.HandleFunc("/", index)
	http.HandleFunc("/index", shell)

	http.HandleFunc("/queries", queries)

	http.HandleFunc("/shell", shell)

	// standard endpoint
	http.HandleFunc("/ws", engine.Ws)

	// endpoints for speed network
	http.HandleFunc("/query", engine.Request)
	http.HandleFunc("/result", engine.Response)

	// for pages under development
	http.HandleFunc("/dev", dev)

	// Create server with a specified address
	srv := &http.Server{
		Addr: ":1111",
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println()
	log.Println("Shutting down zaradb...")

	// Create a context with timeout to allow active requests to complete
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("zaradb forced to shutdown:", err)
	}

	log.Println("zaradb exiting")

}

// render static shell.html file
func queries(w http.ResponseWriter, r *http.Request) {
	// Read the body of the request
	query, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	result := engine.HandleQueries(string(query))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

// redirect to shell page temporary
func dev(w http.ResponseWriter, r *http.Request) {
	f, err := content.ReadFile("static/dev.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(f))
}

// render static shell.html file
func shell(w http.ResponseWriter, r *http.Request) {
	f, err := content.ReadFile("static/shell.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(f))
}

// redirect to shell page temporary
func index(w http.ResponseWriter, r *http.Request) {
	// TOD create index page
	http.Redirect(w, r, "http://localhost:1111/shell", http.StatusSeeOther)
}
