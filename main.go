// Zaradb lite faset document database
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"zaradb/db"
)

var port = db.PORT

func main() {
	// TODO close programe greatfully.

	db := db.Run("test/")
	defer db.Close()
	fmt.Printf("zara run on :%s\n", port)

	http.Handle("/", http.FileServer(http.Dir("web")))

	http.HandleFunc("/shell", shell)

	// standard endpoint
	http.HandleFunc("/ws", Ws)

	// endpoints for speed network
	http.HandleFunc("/query", Resever)
	http.HandleFunc("/result", Sender)

	http.ListenAndServe(":1111", nil)
}

func shell(w http.ResponseWriter, r *http.Request) {
	// Open the index.html file
	f, err := os.Open("web/shell.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the index.html file to the response
	io.Copy(w, f)
}
