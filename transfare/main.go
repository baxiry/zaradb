package main

import (
	"errors"
	"fmt"

	"github.com/melbahja/goph"
)

func main() {
	fmt.Println("vim-go")
}

func transfare() {

	sshcli, err := goph.NewUnknown("root", addr /*"139.162.100.216"*/, goph.Password(h.getPass()))
	if err != nil {
		return errors.New("new connect err" + err.Error())
	}
	defer sshcli.Close()

}
