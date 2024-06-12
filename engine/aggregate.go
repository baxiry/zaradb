package engine

import (
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
		if string(v[len(v)-1]) == "}" {
			continue
		}

		slice = strings.Split(v, siparator)

		res += `"` + slice[0] + `":"` + slice[1] + `",`
	}

	return res[:len(res)-1] + `}`
}

// fields remove or rename fields
func fields(data []string, fields gjson.Result) []string {

	newKey := ""
	oldKey := ""

	toRemove := make([]string, 0)
	for k, v := range fields.Map() {
		fmt.Println(k, v)
		if v.String() == "0" {
			toRemove = append(toRemove, k)
		} else {
			newKey = v.String()
			oldKey = k
		}
	}
	println("toRemove")
	for i := 0; i < len(data); i++ {
		for _, k := range toRemove {
			data[i], _ = sjson.Delete(data[i], k)
		}
	}
	for i := 0; i < len(data); i++ {
		data[i] = reKey(oldKey, newKey, data[i])
	}

	return data
}
