package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/melbahja/goph"
)

type Helper struct{}

var (
	h  Helper
	wg sync.WaitGroup
	mt sync.Mutex

	rp          = getRootPath()
	newHosts    = rp + "new.host"
	activeHosts = rp + "active.host"
	// disactiveHosts = rp + "disactive.host"
	// clientsName    = rp + "clients.name"
	// botLine        = rp + "testBot"
	// statusfile     = rp + "statusb.json"
)

func (Helper) hostIsActive(ctx context.Context, host string) (bool, error) {
	//ch := make(chan bool)
	client, err := goph.NewUnknown("root", host, goph.Password(h.getPass())) //getPass()
	if err != nil {
		//fmt.Println("err with connect ", err)
		return false, err
	}
	defer client.Close()
	select {
	case <-ctx.Done():
		fmt.Println("cancel by context")
		fmt.Println()
		return false, ctx.Err()

	default:
		time.Sleep(time.Millisecond * 100)
		return true, nil
	}
	//return true
}

func (Helper) cleanFile(fname string) {
	err := ioutil.WriteFile(fname, []byte(""), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

// appendData append data to errors.hots file
func (Helper) appendData(newEerror string) {

	file, err := os.OpenFile("errors.host", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Could not open errors.host")
		return
	}
	defer file.Close()

	_, err = file.WriteString(newEerror)

	if err != nil {
		fmt.Println("Could not write text to errors.host")
		return
	}
}

func (Helper) logError(fname, dataError string) {

	//Append second line
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	if _, err := file.WriteString(dataError + "\n"); err != nil {
		log.Fatal(err)
	}

}

// run
func main() {

	hosts, err := h.load(activeHosts)
	if err != nil {
		fmt.Println(err)
		//h.appendData(err.Error() + "\n")
	}
	fmt.Println("lenght of hosts", len(hosts))

	wg.Add(len(hosts))
	for _, host := range hosts {

		ahost := host

		go func(ahost string) {
			defer wg.Done()

			ctx := context.Background()
			cCtx, cancelFunc := context.WithCancel(ctx)

			active, err := h.hostIsActive(cCtx, ahost)

			if active {

				fmt.Println(ahost, "active")
			} else {

				fmt.Println(ahost, "  disactive")
				// mt.Lock(); defer mt.Unlock() // uncomment to avoid data races
				h.logError("error.host", ahost+"  "+err.Error())

				hosts = h.removeItem(ahost, hosts)

			}
			time.Sleep(time.Second * 10)
			cancelFunc()

		}(ahost)
	}
	wg.Wait()

	fmt.Println("lenght of hosts", len(hosts))

	// update active.host file (remove all disactive host)
	h.updateActive(hosts)

	for _, host := range hosts {

		//host := host
		//go func(host string) {
		fmt.Println("download from:", host)

		//err := h.download(host)
		if err != nil {
			//h.appendData(host + "  " + err.Error() + "\n")
			fmt.Println(host, err.Error())
		}
	}
	fmt.Println("sleep")

}

func (Helper) update(file string, list []string) error {
	data := ""
	for _, item := range list {
		data += item + "\n"
	}
	err := os.WriteFile(file, []byte(data), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

// appendAddress appends new address to addressfile
func (Helper) updateActive(hosts []string) {
	f, err := os.Create(activeHosts)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for _, host := range hosts {
		if _, err := f.WriteString(host + "\n"); err != nil {
			log.Println(err)
		}
	}
}

// appendAddress appends new address to addressfile
func (Helper) appendAddr(file, data string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}

// TODO create func that return root path depand on os
func getRootPath() string {

	os := runtime.GOOS
	if os == "linux" {
		return "/root/saverbot/"
	}
	return ""
}

//  zipfile.zip and clientName
func (Helper) remoteZip(cli *goph.Client, dir string) error {
	// zip the client bot app
	// zip -r test.zip testbot
	cmd, err := cli.Command("tar", "-czf", dir+".tar.gz", dir)
	if err != nil {
		return err
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// download zips and downloads remote dir
func (Helper) download(addr string) error {
	sshcli, err := goph.NewUnknown("root", addr /*"139.162.100.216"*/, goph.Password(h.getPass()))
	if err != nil {
		return errors.New("new connect err" + err.Error())
	}
	defer sshcli.Close()

	// get all remote file|dir
	dirs, err := h.ls(addr)
	if err != nil {
		return errors.New("remote ls err: " + err.Error())
	}

	for _, dirbot := range dirs {

		if strings.HasSuffix(dirbot, "-bot") {

			err = h.remoteZip(sshcli, dirbot)
			if err != nil {
				return errors.New("remote zip err: " + err.Error())
			}

			//err = sshcli.Download("/root/"+dirbot+".tar.gz", dirbot+".tar.gz")
			//if err != nil {
			//	return errors.New("download" + err.Error())
			//}
		}
	}

	return nil
}

// to scure app read pass form seprite file
func (Helper) getPass() string {
	data, err := ioutil.ReadFile(".mypass")
	if err != nil {
		return err.Error()
	}
	psw := string(data)
	return psw[:len(psw)-1]
}

// ls list remote file/dir
func (Helper) ls(addr string) ([]string, error) {
	sshcli, err := goph.NewUnknown("root", addr, goph.Password(h.getPass()))
	if err != nil {
		return nil, err
	}
	defer sshcli.Close()

	out, err := sshcli.Run("ls")
	if err != nil {
		return nil, err
	}
	dirs := strings.Split(string(out), "\n")
	return dirs, nil

}

// removeItem remove Item string from list and return new list
func (Helper) removeItem(item string, list []string) []string {
	newList := make([]string, 0)
	for _, v := range list {
		if item != v {
			mt.Lock()
			newList = append(newList, v)
			mt.Unlock()
		}
	}
	return newList
}

// load loads file and return hosts address as []stirng
func (Helper) load(file string) ([]string, error) {

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	hs := strings.Split(string(data), "\n")

	list := make([]string, 0)

	for _, v := range hs {

		h := strings.Replace(v, " ", "", -1)
		if len(h) < 3 {
			continue
		}
		list = append(list, h)
	}

	return h.unique(list), nil
}

// filterList make list unique
func (Helper) unique(list []string) []string {
	mp := make(map[string]bool)
	for _, h := range list {
		mp[h] = true
	}
	ulist := make([]string, 0)
	for h := range mp {
		if h == "" {
			break
		}
		ulist = append(ulist, h)
	}
	return ulist
}

// checkErr check error if exeste and close program
func checkErr(info string, err error) {
	if err != nil {
		log.Println(info, err)
		os.Exit(0)
	}
}
