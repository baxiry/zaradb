package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type IP struct {
	Query string `json:"query"`
}

func getip2() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func main() {
	ch := make(chan bool, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "html")

		fmt.Fprintf(w, `<h3>Host Address : %s<br>Client name : %s<br><br><br>
                       </h3><h2>Worning: type /exit to exit tis application</h2>`,
			getip2(), getClientName("."))

	})

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "will Done") //": %s\n", r.URL.Path)
		ch <- true
	})

	go func() {
		fmt.Println(http.ListenAndServe(":80", nil))
	}()

	go func() {
		if <-ch {
			fmt.Println("Done")
			time.Sleep(time.Second * 2)
			os.Exit(0)
		}
	}()
	for {
		time.Sleep(time.Minute)
	}
}

// get client name
func getClientName(path string) string {

	// Open the directory.
	outputDirRead, err := os.Open(path)
	if err != nil {
		return err.Error()
	}

	// Call Readdir to get all files.
	outputDirFiles, err := outputDirRead.Readdir(0)
	if err != nil {
		return err.Error()
	}

	// Loop over files.
	for _, dir := range outputDirFiles {

		// Print name.
		if dir.IsDir() && strings.HasSuffix(dir.Name(), "-bot") {

			pathFile := dir.Name() + "/" + dir.Name() + ".info"
			fmt.Println("open :", pathFile)

			data, err := ioutil.ReadFile(pathFile)
			if err != nil {
				return err.Error()
			}
			return string(data)
		}
	}
	data, err := ioutil.ReadFile("hamza-bot.info")
	if err != nil {
		return err.Error()
	}
	return string(data)

	//return "no client name"
}
