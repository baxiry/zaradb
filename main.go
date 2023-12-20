// Zaradb lite faset document database
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"zaradb/dblite"
)

func js(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/shell.js")
}

func main() {
	// TODO close programe greatfully.

	db := dblite.Run("test/")
	defer db.Close()
	fmt.Printf("zara run on :%s\n", dblite.PORT)

	http.HandleFunc("/web/shell.js", js)

	http.HandleFunc("/", shell)

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
