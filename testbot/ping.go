package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println(ping("http://localhost/info"))
}

func ping(address string) string {

	req, err := http.Get(address)
	if err != nil {
		fmt.Println("fuck error ", err)
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
