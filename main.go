package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/melbahja/goph"
)

type Host struct {
	Address  string
	Password string
}

func main() {
	hosts, err := loadHosts("hosts.json")
	if err != nil {
		log.Fatalln("err whith loadHosts() function:\n", err)
	}

	client, err := goph.NewUnknown("root", hosts[0].Address, goph.Password(hosts[0].Password))
	checkErr("goph.NewUnknown():", err)
	defer client.Close()
	log.Println("ssh client oppend, Done")
	/*
		// zip the client bot app
		cmd, err := client.Command("zip", "-r", "lilgo.zip ", "lilgo")
		checkErr("client.Command():", err)

		err = cmd.Run()
		checkErr("cmd.Run():", err)
		log.Println("ziped remot file, Done")

		// Download the zeppet bot app
		err = client.Download("/root/lilgo.zip", "lilgo.zip")
		checkErr("err with client.Download()", err)
		log.Println("Download botApp.zip, Done")
	*/
	// Upload new bot app to new host
	err = client.Upload("hosts.json", "/root/hosts.json")
	checkErr("error with deploy(): ", err)

	// run lineBot in new host
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		output, err := client.Run("ls")
		checkErr("client.Run():", err)
		fmt.Println(string(output))
	}()
	wg.Wait()
	//time.Sleep(time.Second * 5)

	// mybe we need enabling bot as a service with systemkd

	// check new client in clientFile one per huor

	// deploying bot to this client

	// Check is evrything is ok by executing ls command
	//out, err := client.Run("ls")

	//if err != nil {
	//	log.Fatal(err)
	//}
}

func loadHosts(file string) ([]Host, error) {

	hosts := make([]Host, 5)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

// checkErr check error if exeste and close program
func checkErr(at string, err error) {
	if err != nil {
		log.Println(at, err)
		os.Exit(0)
	}
}

// zip -r lilgo.zip lilgo
