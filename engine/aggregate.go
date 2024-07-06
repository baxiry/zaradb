package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

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

		if gjson.Parse(slice[1]).Type == 2 {
			res += `"` + slice[0] + `":` + slice[1] + `,`
			continue
		}
		res += `"` + slice[0] + `":"` + slice[1] + `",`
	}

	return res[:len(res)-1] + `}`
}

// fields remove or rename fields
func reFields(data []string, fields gjson.Result) []string {

	newKey := []string{}
	oldKey := []string{}

	toRemove := make([]string, 0)

	for k, v := range fields.Map() {
		if v.String() == "0" {
			toRemove = append(toRemove, k)
		} else {
			newKey = append(newKey, v.String())
			oldKey = append(oldKey, k)
		}
	}

	//remove fields
	for i := 0; i < len(data); i++ {
		for _, k := range toRemove {
			data[i], _ = sjson.Delete(data[i], k)
		}
	}

	//reName fields
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(newKey); j++ {
			data[i] = reKey(oldKey[j], newKey[j], data[i])
		}
	}

	return data
}

func orderBy(param string, reverse int, data []string) (list []string) {

	objects := []gjson.Result{}
	for _, v := range data {
		objects = append(objects, gjson.Parse(v))
	}
	// sort here

	// check type
	typ := objects[0].Get(param).Type
	fmt.Println("type is ", typ)

	if typ == 2 {
		list = sortNumber(param, reverse, objects)
	}
	if typ == 3 {
		list = sortString(param, reverse, objects)
	}

	return list
}

func sortNumber(field string, reverse int, list []gjson.Result) []string {
	max := len(list)
	var tmp gjson.Result

	element := list[0]
	for max != 1 {
		for i := 1; i < max; i++ {
			if element.Get(field).Num < list[i].Get(field).Num {
				tmp = list[i]
				list[i] = element
				element = tmp
			}

			if i == max-1 {
				tmp = list[i]
				list[i] = element
				element = tmp
			}
		}
		max--
	}

	list[0] = element
	res := []string{}
	if reverse != 1 {
		for i := 0; i < len(list); i++ {
			res = append(res, list[i].String())
			fmt.Println(list[i].String())
		}
		return res
	}

	for i := len(list) - 1; i >= 0; i-- {
		fmt.Println(list[i].String())
		res = append(res, list[i].String())
	}

	return res
}

// TODO  consider specific type.
func sortString(field string, reverse int, list []gjson.Result) []string {
	max := len(list)
	var tmp gjson.Result

	element := list[0]
	for max != 1 {
		for i := 1; i < max; i++ {
			if element.Get(field).Str < list[i].Get(field).Str {
				tmp = list[i]
				list[i] = element
				element = tmp
			}

			if i == max-1 {
				tmp = list[i]
				list[i] = element
				element = tmp
			}
		}
		max--
	}

	list[0] = element

	res := []string{}
	if reverse != 1 {
		for i := 0; i < len(list); i++ {
			res = append(res, list[i].String())
			fmt.Println(list[i].String())
		}
		return res
	}

	for i := len(list) - 1; i >= 0; i-- {
		res = append(res, list[i].String())
		fmt.Println(list[i].String())
	}
	return res
}
