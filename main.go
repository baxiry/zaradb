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

// main
func main() {
	// TODO: Close program gracefully.

	db := engine.NewDB("test.db")
	db.CreateCollection("test")
	defer db.Close()

	fmt.Printf("interacte with zaradb through %s:%s\n", Host, Port)

	http.Handle("/static/", http.FileServer(http.FS(content)))

	http.HandleFunc("/", index)
	http.HandleFunc("/index", shell)

	http.HandleFunc("/shell", shell)

	// standard endpoint
	http.HandleFunc("/ws", engine.Ws)

	// endpoints for speed network
	http.HandleFunc("/query", engine.Request)
	http.HandleFunc("/result", engine.Response)

	// for pages under development
	http.HandleFunc("/dev", dev)

	log.Println(http.ListenAndServe(":1111", nil))
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
	// TODO create index page
	http.Redirect(w, r, "http://localhost:1111/shell", http.StatusSeeOther)
}
