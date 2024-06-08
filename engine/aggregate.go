package engine

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// TODO implemente replace keys

const siparator = "_:_"

// reKey renames json feild
func reKey(oldkey, newkey, json string) string {

	list := []string{}
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {

		if key.String() == oldkey {
			list = append(list, newkey+siparator+value.String())
		} else {

			list = append(list, key.String()+siparator+value.String())
		}

		return true
	})
	slice := []string{}

	res := "{"
	for _, v := range list {
		slice = strings.Split(v, siparator)
		res += `"` + slice[0] + `":"` + slice[1] + `",`
	}

	return res[:len(res)-1] + `}`
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

// TODO this func for studing
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
