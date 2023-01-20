package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	dbFile := "../example.db"

	file, err := Opendbs(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data := "tested data ok"
	AppendData(file, data)

	getedData := GetVal(file, 0, 14)

	fmt.Println(getedData)
}

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return strings.Split(args[1], ".")
}
