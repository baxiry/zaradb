package main

import (
	"fmt"

	"github.com/sfreiberg/simplessh"
)

func main() {
	/*
		Leave privKeyPath empty to use $HOME/.ssh/id_rsa.
		If username is blank simplessh will attempt to use the current user.
	*/
	client, err := simplessh.ConnectWithKeyFile("158.247.195.235", "root", "/root/.ssh/id_rsa")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	output, err := client.Exec("uptime")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Uptime: %s\n", output)
}
