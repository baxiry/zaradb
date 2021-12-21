package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

var (
	cmd  = "scp"
	args = "root@139.162.121.240:/root/.bashrc ./"
	pass = "d7ombot123"
)

func main() {

	cmd := exec.Command(cmd, args)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, pass)
	}()

	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)

}
