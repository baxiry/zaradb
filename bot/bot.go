package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
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

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "html")

		out := execute("pwd")

		fmt.Fprintf(w, "<h3>My Ip Address : %s<br>My Dir work : %s<br> Worning: type /exit to exit tis application</h3>", getip2(), out)
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

func execute(cmd string) string {
	out, err := exec.Command(cmd).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}
	return string(out)

}
