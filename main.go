package main

import (
	"fmt"
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
	if err != nil {
		fmt.Printf("%T\n", err.Error())
		os.Exit(1)
	}
	defer client.Close()

	// zip the client bot app
	cmd, err := client.Command("zip", "-r", "/root/lilgo.zip ", "/root/lilgo")
	if err != nil {
		fmt.Println("cmd error:", err)
	}
	cmd.Run()
	if err != nil {
		fmt.Println("err with cmd.Run()", err)
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("err with cmd.Output()", err)
	}
	fmt.Println("output: ", string(output))

	// Download the zeppet bot app
	//err = client.Download("/root/lilgo.zip", "lilgo.zip")
	//fmt.Println("error is :", err)

	// Upload new bot app to new host

	// run lineBot in new host

	// mybe we need enabling bot as a service with systemd

	// check new client in clientFile one per huor

	// deploying bot to this client

	// Check is evrything is ok by executing ls command
	//out, err := client.Run("ls")

	//if err != nil {
	//	log.Fatal(err)
	//}

	// Get your output as []byte.
	fmt.Println("list is :\n", string(output))
}
