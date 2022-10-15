package main

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

func getJson() {

	d, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println(err)
	}

	value := gjson.Get(string(d), "data._id")
	println(value.String())
	value = gjson.Get(string(d), "data._idd")
	println(value.String())

}
