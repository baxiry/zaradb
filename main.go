package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/melbahja/goph"
)

type Helper struct{}

var h Helper

func main() {

}

// writeData updates/rewrites data into file
func (Helper) writeData(file, data string) error {
	err := os.WriteFile(file, []byte(data+"\n"), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

// loadDisactive load addresses of disactive hosts
func (Helper) disactiveHosts(path string) ([]string, error) {

	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(bin), "\n"), nil
}

// appendAddress appends new address to addressfile
func (Helper) appendAddress(file, data string) {
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
func (Helper) unique(data []string) []string {
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

type Bot struct {
	Owner   string `json:"owner"`
	Address string `json:"address"`
	Active  bool   `json:"active"`
}

// activeHosts filter hosts and return just active hostes
func (Helper) activeHosts(bots []Bot) []Bot {
	activeBots := make([]Bot, 0)
	for _, bot := range bots {
		if bot.Active {
			activeBots = append(activeBots, bot)
		} else {
			h.appendAddress("disactive.json", bot.Address)
		}
	}
	return activeBots
}

// return list of bots type
func (Helper) loadBots(file string) ([]Bot, error) {

	bots := make([]Bot, 5)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &bots)
	if err != nil {
		return nil, err
	}
	return bots, nil
}

// TODO test zip function
//  zipfile.zip and clientName
func (Helper) zip(client *goph.Client, outfile, dir string) error {
	// zip the client bot app
	cmd, err := client.Command("zip", "-r", outfile+".zip", dir)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// exitBot
func (Helper) exitBot() {
	ch := make(chan bool, 1)
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
			os.Exit(0)
		}
	}()
}

// send exitbot
func (Helper) sendExit(address string) {
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

// checkErr check error if exeste and close program
func checkErr(at string, err error) {
	if err != nil {
		log.Println(at, err)
		os.Exit(0)
	}
}
