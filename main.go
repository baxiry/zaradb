// Zaradb lite faset document database
package main

import (
	"embed"
	"fmt"
	"net/http"
<<<<<<< HEAD
	"zaradb/db"
)

//go:embed static/*
var staticDir embed.FS

var port = db.PORT
=======
	"zaradb/database"
)

//go:embed  static
var content embed.FS
>>>>>>> 4b01cc4 (adompting sqlite & and reorginazing zara project)

func main() {
	// TODO close programe greatfully.

	db, err := database.NewDB("test.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
<<<<<<< HEAD
	fmt.Printf("zara run on :%s\n", "localhost:"+port+"/shell")

	http.Handle("/static/", http.FileServer(http.FS(staticDir)))
=======

	fmt.Printf("zara run on :%s\n", database.PORT)

	//http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(content)))) // http.FileServer(http.Dir("static"))))
	http.Handle("/", http.FileServer(http.FS(content)))
>>>>>>> 4b01cc4 (adompting sqlite & and reorginazing zara project)

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
<<<<<<< HEAD

	http.ServeFile(w, r, "shell.html")
=======
	// Open the index.html file
	f, err := content.ReadFile("static/shell.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy the index.html file to the response
	fmt.Fprint(w, string(f))
>>>>>>> 4b01cc4 (adompting sqlite & and reorginazing zara project)
}

/*

func ServeStaticFiles(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/static/")
		fmt.Println("file path is : ", path)
		// Remove "/static/" prefix

		// Open the file
		f, err := root.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer f.Close()

		// Determine MIME type based on extension
		contentType := mime.TypeByExtension(path)
		if contentType == "" {
			// Fallback for unknown extensions
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)

		// Serve the file
		http.ServeFile(w, r, path)
	})
}


*/
