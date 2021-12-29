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

var wg sync.WaitGroup

type Host struct {
	Address    string `json:"address"`
	Password   string `json:"password"`
	ClientName string `json:"clientName"`
	Active     bool   `json:"active"`
}

func main() {
	hosts, err := loadHosts("hosts.json")
	if err != nil {
		log.Fatalln("err whith loadHosts() function:\n", err)
	}

	allhosts, err := loadHosts("hosts.json")
	checkErr("loadHosts:", err)

	activeHosts := filterActive(allhosts)
	fmt.Println("all hosts is : ")
	for _, host := range allhosts {
		fmt.Println(host.ClientName)
	}

	fmt.Println()
	fmt.Println("active hosts is : ")
	for _, host := range activeHosts {
		fmt.Println(host.ClientName)
	}
	os.Rename("test", activeHosts[0].ClientName)

	os.Exit(0)

	// importent
	for _, host := range hosts {
		host := host

		fmt.Println(host.Address)
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := goph.NewUnknown("root", host.Address, goph.Password(host.Password))
			checkErr("goph.NewUnknown():", err)
			defer client.Close()
			log.Println("ssh client oppend, Done")
			// Upload new bot app to new host
			err = client.Upload("hosts.json", "/root/hosts.json")
			checkErr("error with deploy(): ", err)

			// run lineBot in new host

			output, err := client.Run("hostname -I")
			checkErr("client.Run():", err)
			fmt.Println("remote ip address is : ", string(output))

		}()
	}
	wg.Wait()

	// mybe we need enabling bot as a service with systemkd

	// check new client in clientFile one per huor

	// deploying bot to this client
}

func filterActive(hosts []Host) []Host {
	activeHosts := make([]Host, 0)
	for _, host := range hosts {
		if host.Active {
			activeHosts = append(activeHosts, host)
		}
	}
	return activeHosts
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

// zip -r lilgo.zip lilgo
// Check is evrything is ok by executing ls command
//out, err := client.Run("ls")

//if err != nil {
//	log.Fatal(err)
//}
