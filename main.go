package main

import (
	"fmt"
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
