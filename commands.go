package main

import (
	"fmt"
	"os"
	"strings"
)

// cli functions

const hints = `tap helpe to get help massage`

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return strings.Split(args[1], ".")
}
