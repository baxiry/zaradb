// Zaradb lite faset document database
package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"zaradb/engine"
)

//go:embed  static
var content embed.FS

func main() {
	// TODO: Close program gracefully.

	db := engine.NewDB("test.db")
	db.CreateCollection("test")
	defer db.Close()

	fmt.Printf("interacte with zaradb through %s:%s\n", Host, Port)

	http.Handle("/static/", http.FileServer(http.FS(content)))

	http.HandleFunc("/shell", shell)

	// not importent
	http.HandleFunc("/dev", dev)

	// standard endpoint
	http.HandleFunc("/ws", engine.Ws)

	// endpoints for speed network
	http.HandleFunc("/query", engine.Request)
	http.HandleFunc("/result", engine.Response)

	log.Println(http.ListenAndServe(":1111", nil))
}

// redirect to shell page temporary
func dev(w http.ResponseWriter, r *http.Request) {
	// TODO create index page
	//http.Redirect(w, r, "http://localhost:1111/shell", http.StatusSeeOther)
	f, err := content.ReadFile("static/dev.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(f))
}

// render static shell.html file
func shell(w http.ResponseWriter, r *http.Request) {
	// Open the index.html file
	f, err := content.ReadFile("static/shell.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the index.html file to the response
	fmt.Fprint(w, string(f))
}
