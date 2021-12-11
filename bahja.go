package main

import (
	"fmt"
	"log"

	"github.com/melbahja/goph"
)

func main() {

	client, err := goph.New("root", "158.247.195.235", goph.Password("e1J=xHytspxhhscx"))
	if err != nil {
		log.Fatal(err)
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	out, err := client.Run("ls go")

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println("go dir is :\n", string(out))

	// Execute your command.
	out, err = client.Run("ls .go")

	if err != nil {
		log.Fatal(err)
	}

	// Get your output as []byte.
	fmt.Println("dot go dir is :\n", string(out))

}
