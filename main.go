package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/melbahja/goph"
)

func main() {

	allhosts, err := loadHosts("hosts.json")
	activeHosts := filterActive(allhosts)
	checkErr("loadHosts:", err)

	client, err := goph.NewUnknown("root", activeHosts[0].Address, goph.Password(psw))
	checkErr("goph.NewUnknown():", err)
	defer client.Close()
	output, err := client.Run("ls")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))

	//zip(client, activeHosts[0].ClientName, activeHosts[0].ClientName)

	os.Exit(0)

	os.Rename("test", activeHosts[0].ClientName)

	// importent
	for _, host := range activeHosts {
		host := host

		fmt.Println(host.Address)
		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := goph.NewUnknown("root", host.Address, goph.Password(psw))
			checkErr("goph.NewUnknown():", err)
			defer client.Close()
			log.Println("ssh client oppend, Done")
			// Upload new bot app to new host
			err = client.Upload("disactive.json", "/root/hosts.json")
			checkErr("error with deploy(): ", err)

			// run lineBot in new host

			output, err := client.Run("hostname -I")
			checkErr("client.Run():", err)
			fmt.Println("remote ip address is : ", string(output))

		}()
	}

	// check new client in clientFile one per huor

	// deploying bot to this client
	wg.Wait()
}

var wg sync.WaitGroup

type Helper struct{}

var h Helper

// writeData update/rewrite data into file
func (Helper) writeData(file, data string) error {
	err := os.WriteFile(file, []byte(data+"\n"), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
	//defer f.Close()
	//if _, err := f.WriteString(data + "\n"); err != nil {
	//	log.Println(err)
	//}
}

// loadDisactive load addresses of disactive hosts
func (Helper) loadDisactive(path string) ([]string, error) {

	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bin), "\n"), nil
}

// appendData append new address to addressfile
func (Helper) appendData(file, data string) {
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}

// filterList make list unique
func (Helper) filterList(data []string) []string {
	mp := make(map[string]bool)
	for _, h := range data {
		mp[h] = true
	}
	hosts := make([]string, 0)
	for h := range mp {
		if h == "" {
			break
		}
		hosts = append(hosts, h)
	}
	return hosts
}

type Host struct {
	Address    string `json:"address"`
	ClientName string `json:"clientName"`
	Active     bool   `json:"active"`
}

func filterActive(hosts []Host) []Host {
	activeHosts := make([]Host, 0)
	for _, host := range hosts {
		if host.Active {
			activeHosts = append(activeHosts, host)
		} else {
			h.appendData("disactive.json", host.Address)
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

// TODO test zip function
//  zipfile.zip and clientName
func zip(client *goph.Client, zipfile, dir string) error {
	// zip the client bot app
	cmd, err := client.Command("zip", "-r", zipfile+".zip", dir)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

/*
	// zip the client bot app
	cmd, err := client.Command("zip", "-r", "lilgo.zip", "lilgo")
	checkErr("client.Command():", err)

	err = cmd.Run()
	checkErr("cmd.Run():", err)
	log.Println("ziped remot file, Done")

	// Download the zeppet bot app
	err = client.Download("/root/lilgo.zip", "lilgo.zip")
	checkErr("err with client.Download()", err)
	log.Println("Download botApp.zip, Done")
*/

/*
zip -r lilgo.zip lilgo
// Check is evrything is ok by executing ls command
out, err := client.Run("ls")

if err != nil {
	log.Fatal(err)
}
*/

/*
exitbot("http://localhost:8080/" + "exit")
	fmt.Println("done exit bot")
*/

// exitBot
func exitBot() {
	ch := make(chan bool, 1)
	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "will Done") //": %s\n", r.URL.Path)
		ch <- true
	})
	go func() {
		fmt.Println(http.ListenAndServe(":8080", nil))
	}()

	go func() {
		if <-ch {
			fmt.Println("Done")
			os.Exit(0)
		}
	}()
}

// send exitbot
func sendExit(address string) {
	resp, err := http.Get("http://" + address + ":8080/exit")
	if err != nil {
		log.Fatal("Error getting response. ", err)
	}
	defer resp.Body.Close()

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}

	fmt.Printf("body is : %s\n", body)
}

const psw = "d7ombot123"
