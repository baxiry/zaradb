package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// not implemented yet
func aggrigate(query gjson.Result) string {
	data, err := getData(query)
	if err != nil {
		return err.Error()
	}

	group := query.Get("group")
	if _id := group.Get("_id"); !_id.Exists() {
		return "a group specification must include an _id"
	}

	mapData := map[string]string{}
	message := ""

	_id := group.Get("_id").Str

	// TODO parse data and exclude invalide objects
	if gjson.Get(data[0], _id).Str == "" {
		return "field '" + _id + "' is not exists"
	}

	group.ForEach(func(key, val gjson.Result) bool {
		if key.Str != "_id" && val.Type != 5 {
			message = "The field '" + key.Str + "' must be an accumulator object!"
			return false
		}

		switch val.Type {
		case 3:
			for _, obj := range data {
				field := gjson.Get(obj, val.String()).String()
				json, _ := sjson.Set("", key.String(), field)
				mapData[field] = json
			}

		case 5:
			val.ForEach(func(opr, fld gjson.Result) bool { // opperation & field name

				switch opr.Str {
				case "$count":
					counted := count(fld.Str, data)
					for _id, count := range counted {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, count)
					}

				case "$max":
					maxs, err := max(_id, fld, data)
					if err != nil {
						message = err.Error()
						return false
					}
					for _id, max := range maxs {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, max)
					}

				case "$min":
					mins, err := min(_id, fld, data)
					if err != nil {
						message = err.Error()
						return false
					}

					for _id, min := range mins {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, min)
					}

				case "$sum":
					sums := sum(_id, fld.Str, data)
					for _id, sumd := range sums {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, sumd)
					}

				case "$avg":
					avrs := average(_id, fld.Str, data)
					for _id, avr := range avrs {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, avr)
					}

				default:
					message = "unknown '" + opr.Str + "' aggrigate operator !"
				}

				return true
			})

		default:
			message = "opss! something wrrong!"
		}

		return true
	})

	result := "["
	for _, val := range mapData {
		result += val + ","
	}
	if message != "" {
		return message
	}
	return result[:len(result)-1] + "]"
}

// gets mines vlaues per _id (e.g. name)
func min(_id string, field gjson.Result, records []string) (mp map[string]float64, err error) {
	min := float64(9223372036854775807) // we need max float
	mp = map[string]float64{}

	// init mp by max float val
	for _, record := range records {
		id := gjson.Get(record, _id).Str
		mp[id] = min
	}

	switch field.Type {
	case 3:
		for _, record := range records {
			id := gjson.Get(record, _id).Str        // name of record
			val := gjson.Get(record, field.Str).Num // value of compared field
			if mp[id] > val {
				mp[id] = val
			}
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":

				for _, record := range records {

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 * arg2
					if mp[id] > val {
						mp[id] = val
					}
				}
			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 + arg2
					if mp[id] > val {
						mp[id] = val
					}
				}

			case "$sub":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 - arg2
					if mp[id] > val {
						mp[id] = val
					}
				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field
					val := arg1 / arg2

					if mp[id] > val {
						mp[id] = val
					}
				}

			default:
				err = fmt.Errorf("unknown %s", op)
			}

			return false
		})

		fmt.Println(field, "is operation")
		fmt.Println()
		fmt.Println(field.Get("$min"))
	}

	return mp, nil
}

func sum(_id, field string, records []string) (mp map[string]float64) {
	mp = map[string]float64{}
	for _, record := range records {
		id := gjson.Get(record, _id).Str    // name of record
		val := gjson.Get(record, field).Num // value of sumd field
		mp[id] += val
	}
	return mp
}

// max gets vlaues per _id (e.g. name)
func max(_id string, field gjson.Result, records []string) (mp map[string]float64, err error) {
	max := float64(-9223372036854775808) // we need min floa
	mp = map[string]float64{}

	// init mp by mines float val
	for _, record := range records {
		id := gjson.Get(record, _id).Str
		mp[id] = max
	}
	/*
		for _, record := range records {
			id := gjson.Get(record, _id).Str        // name of record
			val := gjson.Get(record, field.Str).Num // value of compared field
			if mp[id] < val {
				mp[id] = val
			}
		}
	*/
	switch field.Type {
	case 3:
		for _, record := range records {
			id := gjson.Get(record, _id).Str        // name of record
			val := gjson.Get(record, field.Str).Num // value of compared field
			if mp[id] < val {
				mp[id] = val
			}
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":

				for _, record := range records {

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 * arg2
					if mp[id] < val {
						mp[id] = val
					}
				}
			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 + arg2
					if mp[id] < val {
						mp[id] = val
					}
				}

			case "$sub":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field

					val := arg1 - arg2
					if mp[id] < val {
						mp[id] = val
					}
				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str                 // name of record
					arg1 := gjson.Get(record, args.Get("0").Str).Num // value of compared field
					arg2 := gjson.Get(record, args.Get("1").Str).Num // value of compared field
					val := arg1 / arg2

					if mp[id] < val {
						mp[id] = val
					}
				}

			default:
				err = fmt.Errorf("unknown %s operator", op)
			}

			return false
		})

		fmt.Println(field, "is operation")
		fmt.Println()
		fmt.Println(field.Get("$min"))
	}
	return mp, nil
}

// not implemented yet
func average(_id, field string, records []string) (mp map[string]float64) {
	mp = map[string]float64{}

	fieldCount := make(map[string]float64)

	for _, record := range records {
		id := gjson.Get(record, _id).Str // name of record
		val := gjson.Get(record, field)  // value of sumd field

		//
		if val.Exists() {
			fieldCount[id]++
		}

		mp[id] += val.Num

	}
	for fld, count := range fieldCount {
		mp[fld] = mp[fld] / count
	}

	return mp
}

func count(field string, records []string) (mp map[string]int) {

	mp = map[string]int{}
	for _, record := range records {
		mp[gjson.Get(record, field).String()]++
	}
	fmt.Println("result:  ", mp)
	return mp
}

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
