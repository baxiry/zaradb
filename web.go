package main

import (
	"fmt"
	"net/http"
)

func mainc() {
	fmt.Println("service is runing")
	http.HandleFunc("/", indexPage)
	http.ListenAndServe(":8080", nil)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "service is still runing")
}
