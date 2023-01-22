package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	url = "http://158.247.195.235:8001"
)

func readSerial() string {
	data, err := ioutil.ReadFile("./serial.txt")
	if err != nil {
		fmt.Println("you don't have serial file")
		os.Exit(0)
	}
	serial := strings.Replace(string(data), "\n", "", 1)
	return serial
}
func main() {

	serial := readSerial()
	userAuth(serial)
	expiration(serial)
}

func expiration(serial string) {
	resp, err := http.Get(url + "/time?serial=" + serial)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	btime, _ := ioutil.ReadAll(resp.Body)
	day, err := strconv.Atoi(string(btime))
	if int(day) < 0 {
		fmt.Println("License time has expired")
		os.Exit(0)
	}

	fmt.Printf("%s days left until the license expires\n", string(btime))
}

// check if bot is signup by serial
func userAuth(serial string) {
	resp, err := http.Get(url + "/info?serial=" + serial)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	user, _ := ioutil.ReadAll(resp.Body)
	if len(user) < 5 {
		fmt.Println("no auth ")
		os.Exit(0)
	}
	println(string(user))
}
