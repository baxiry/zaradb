package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var str = fmt.Sprint

func transaction(query gjson.Result) string {

	actions := query.Get("transaction").Array()
	start := "t " + str(len(actions)) + "\n"
	_ = start
	for k, v := range actions {

		fmt.Println(k, v)
	}
	return "not implement yet"
}
