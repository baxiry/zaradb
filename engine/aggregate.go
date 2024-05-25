package engine

import (
	"bytes"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// TODO implemente replace keys

// ReplaceKeys replaces json keys in data using the values in keymap
func ReplaceKeys(data []byte, keymap map[string]string) []byte {
	for kafkaKey, esKey := range keymap {
		old := fmt.Sprintf("\"%s\":", kafkaKey)
		newd := fmt.Sprintf("\"%s\":", esKey)
		data = bytes.Replace(data, []byte(old), []byte(newd), 1)
	}
	// note emplemeted yet
	return data
}

// fields remove or rename fields
func fields(data []string, fields gjson.Result) []string {

	fmt.Println("fields: ")
	toRemove := make([]string, 0)
	for k, v := range fields.Map() {
		fmt.Println(k, v)
		if v.String() == "0" {
			toRemove = append(toRemove, k)
		}
	}
	println("toRemove")
	for i := 0; i < len(data); i++ {
		for _, k := range toRemove {
			data[i], _ = sjson.Delete(data[i], k)
		}
	}

	return data
}
