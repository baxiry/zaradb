package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var str = fmt.Sprint

func transaction(query string) string {

	actions := gjson.Get(query, "transaction").Array()
	start := "t " + str(len(actions)) + "\n"
	_ = start
	for k, v := range actions {

		fmt.Println(k, v)
	}
	return "actions done"
}
