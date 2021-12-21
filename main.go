package main

import (
	"log"
	"os"

	"github.com/melbahja/goph"
)

// zip -r lilgo.zip lilgo

var (
	//address = "172.105.241.152"
	address  = "158.247.195.235"
	password = "e1J=xHytspxhhscx"

	//namefile = "lilgo.zip"
	//path     = "/root/"
)

func main() {

	//client, err := goph.New("root", address, goph.Password(password))
	client, err := goph.NewUnknown("root", address, goph.Password(password))
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

		// Upload new bot app to new host
		err = client.Upload("web.go", "/root/web.go")
		checkErr("error with deploy(): ", err)
	*/
	// run lineBot in new host
	cmd, err := client.Command("/root/web &")
	checkErr("client.Command():", err)

	err = cmd.Run()
	checkErr("cmd.Run():", err)
	log.Println("run bot... Done")

	// mybe we need enabling bot as a service with systemkd

	// check new client in clientFile one per huor

	// deploying bot to this client

	// Check is evrything is ok by executing ls command
	//out, err := client.Run("ls")

	//if err != nil {
	//	log.Fatal(err)
	//}
}

// checkErr check error if exeste and close program
func checkErr(at string, err error) {
	if err != nil {
		log.Println(at, err)
		os.Exit(0)
	}
}
