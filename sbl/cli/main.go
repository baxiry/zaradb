// tool to manage lisences
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var url = "http://158.247.195.235:8001"

func main() {

	var cmd, ip, serial, name string
	flag.StringVar(&cmd, "cmd", "", "The 'cmd' must be 'new', 'update', 'get' to create update or get bot")
	flag.StringVar(&name, "name", "", "The 'name' for name boot ")
	flag.StringVar(&serial, "ser", "", "The 'ser' must be serial")
	flag.StringVar(&ip, "ip", "ip", "The 'ip' must be ip addres")

	// parse flags from command line
	flag.Parse()

	switch cmd {
	case "new":
		createBoot(name, ip) // test this
	case "update":
		update(name, ip, serial)
	case "get":
		info(serial)
		println("get boot info")
	case "delete":
		println("delete boot")
	default:
		fmt.Println("type -help for help message")
	}

	os.Exit(0)
}

// update update boot name or ip or both
func info(serial string) {
	//
	resp, err := http.Get(url + "/info?serial=" + serial)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("rese: ", string(body))
}

// update update boot name or ip or both
func update(name, ip, serial string) {
	if len(serial) == 0 {
		fmt.Println("you messing serial")
		return
	}
	//
	resp, err := http.Get(url + "/update?serial=" + serial + "&name=" + name + "&ip=" + ip)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("rese: ", string(body))
}

func createBoot(name, ip string) {
	//
	resp, err := http.Get(url + "/new?serial=" + newSerial() + "&name=" + name + "&ip=" + ip)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("rese: ", string(body))
}
func newSerial() (serial string) {
	chars := []string{
		"q", "w", "e", "r", "t", "y", "u", "i", "o",
		"a", "s", "d", "f", "g", "h", "l", "k", "j",
		"A", "B", "C", "D", "E", "F", "J", "H", "E"}

	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 10; i++ {
		serial += chars[rand.Intn(len(chars)-1)]
	}
	return serial
}

/*
func PingDB(db *sql.DB) {
	err := db.Ping()
	panic(err)
}

// initialaze database
func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@/licenses")
	if err != nil {
		panic(err)
	}
	return db
}

*/
