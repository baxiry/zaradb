package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	// root@158.247.195.235      e1J=xHytspxhhsc

	// ssh refers to the custom package above
	conn, err := ssh.Conn("158.247.195.235", "root", "e1J=xHytspxhhsc")
	if err != nil {
		log.Fatal(err)
	}

	output, err := conn.SendCommands("sleep 3", "echo Hello!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))

}
