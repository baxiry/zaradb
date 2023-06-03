package main

import (
	"log"
	"os/user"
)

func userDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
