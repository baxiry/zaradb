// Zaradb lite faset document database
package main

import (
	"embed"
	"fmt"
	"net/http"
	"zaradb/db"
)

//go:embed static/*
var staticDir embed.FS

var port = db.PORT

func main() {
	// TODO close programe greatfully.

	db := db.Run("test/")
	defer db.Close()
	fmt.Printf("zara run on :%s\n", "localhost:"+port+"/shell")

	http.Handle("/static/", http.FileServer(http.FS(staticDir)))

	http.HandleFunc("/shell", shell)

	// standard endpoint
	http.HandleFunc("/ws", Ws)

	// endpoints for speed network
	http.HandleFunc("/query", Resever)
	http.HandleFunc("/result", Sender)

	http.ListenAndServe(":1111", nil)
}

// render static shell.html file
func shell(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "shell.html")
}
