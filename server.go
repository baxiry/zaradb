package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func process(w http.ResponseWriter, r *http.Request) {

	filename := r.URL.Query().Get("filename")
	fmt.Println("filename is : ", filename)

	switch r.Method {
	case "GET":
		data := readFile(filename)
		fmt.Fprintf(w, data)

	case "POST":
		data, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%#v\n", string(data))
		fmt.Println("filename in switch", filename)
		writeFile(filename, string(data)+"\n")

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {

	http.HandleFunc("/saveme", process)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readFile(path string) string {
	data := ""
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data += scanner.Text() + "\n"

		fmt.Println(scanner.Text())
		// do somethin with this data
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}

func writeFile(path, data string) {
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}
