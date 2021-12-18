package main

import (
	"fmt"
	"log"
	"os"

	"github.com/melbahja/goph"
)

var (
	//address = "172.105.241.152"
	address  = "158.247.195.235"
	password = "e1J=xHytspxhhscx"
)

func main() {

	//client, err := goph.New("root", address, goph.Password(password))
	client, err := goph.NewUnknown("root", address, goph.Password(password))
	if err != nil {
		fmt.Printf("%T\n", err.Error())
		os.Exit(1)
	}

	// Defer closing the network connection.
	defer client.Close()

	fmt.Println("close client")

	// Execute your command.
	out, err := client.Run("ls")

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println("list is :\n", string(out))

	// Execute your command.
	out, err = client.Run("ls /")

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println("root list is :\n", string(out))
}
