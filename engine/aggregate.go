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
		return `{"code":0, "status":"a group specification must include an _id"}`
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
					sums, err := sum(_id, fld, data)
					if err != nil {
						message = err.Error()
						return false
					}

					for _id, sumd := range sums {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, sumd)
					}

				case "$avg":
					avrs, err := average(_id, fld, data)
					if err != nil {
						message = err.Error()
						return false
					}
					for _id, avr := range avrs {
						mapData[_id], _ = sjson.Set(mapData[_id], key.Str, avr)
					}

				default:
					message = "unknown '" + opr.Str + "' aggrigate operator !"
					return false
				}

				return true
			})

		default:
			message = "opss! something wrrong!"
			return false
		}

		return true
	})

	if message != "" {
		return message
	}

	filter := query.Get(gmatch)
	limit := query.Get(glimit).Int()
	skip := query.Get(gskip).Int()

	if limit == 0 {
		// what is default setting ?
		limit = 1000
	}

	listdata := []string{}
	for _, val := range mapData {
		if ok, _ := match(filter, val); ok {
			listdata = append(listdata, val)
		}
	}

	// TODO sort listdata here

	result := "["
	for _, val := range listdata {

		if limit == 0 {
			break
		}
		if skip != 0 {
			skip--
			continue
		}

		listdata = append(listdata, val)
		result += val + ","
		limit--
	}

	ln := len(result)
	if ln == 1 {
		return "[]"
	}

	return result[:ln-1] + "]"

}

func sum(_id string, field gjson.Result, records []string) (mp map[string]float64, err error) {
	mp = map[string]float64{}

	switch field.Type {
	case 3:
		for _, record := range records {
			id := gjson.Get(record, _id).Str
			val := gjson.Get(record, field.Str).Num

			mp[id] += val
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":

				for _, record := range records {
					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					mp[id] += arg1 * arg2
				}
			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					mp[id] += arg1 + arg2
				}

			case "$sub":

				for _, record := range records {

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					mp[id] += arg1 - arg2
				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					mp[id] += arg1 / arg2
				}

			default:
				err = fmt.Errorf("unknown %s operator", op)
			}

			return false
		})
		if err != nil {
			return nil, err
		}

		fmt.Println(" we are here ")
	}

	return mp, nil
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
			id := gjson.Get(record, _id).Str
			val := gjson.Get(record, field.Str).Num
			if mp[id] > val {
				mp[id] = val
			}
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":
				for _, record := range records {

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 * arg2
					if mp[id] > val {
						mp[id] = val
					}
				}
			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 + arg2
					if mp[id] > val {
						mp[id] = val
					}
				}

			case "$sub":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 - arg2
					if mp[id] > val {
						mp[id] = val
					}
				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num
					val := arg1 / arg2

					if mp[id] > val {
						mp[id] = val
					}
				}

			default:
				err = fmt.Errorf("unknown %s operator", op)
			}

			return false
		})
		if err != nil {
			return nil, err
		}
	}

	return mp, nil
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

	switch field.Type {
	case 3:
		for _, record := range records {
			id := gjson.Get(record, _id).Str
			val := gjson.Get(record, field.Str).Num
			if mp[id] < val {
				mp[id] = val
			}
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":

				for _, record := range records {

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 * arg2
					if mp[id] < val {
						mp[id] = val
					}
				}
			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 + arg2
					if mp[id] < val {
						mp[id] = val
					}
				}

			case "$sub":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 - arg2
					if mp[id] < val {
						mp[id] = val
					}
				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num
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

		if err != nil {
			return nil, err
		}
	}

	return mp, nil
}

// not implemented yet
func average(_id string, field gjson.Result, records []string) (mp map[string]float64, err error) {
	mp = map[string]float64{}

	fieldCount := make(map[string]float64)

	switch field.Type {
	case 3:
		for _, record := range records {
			id := gjson.Get(record, _id).Str
			val := gjson.Get(record, field.Str).Num

			//if val.Exists() {}
			fieldCount[id]++
			mp[id] += val //.Num
		}

	case 5:

		field.ForEach(func(op, args gjson.Result) bool {
			switch op.Str {
			case "$multiply":

				for _, record := range records {

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					fieldCount[id]++
					mp[id] += arg1 * arg2
				}

			case "$add":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					fieldCount[id]++
					mp[id] += arg1 + arg2
				}

			case "$sub":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num

					val := arg1 - arg2
					mp[id] += val //.Num

					fieldCount[id]++

				}

			case "$div":

				for _, record := range records {
					fmt.Println(record)

					id := gjson.Get(record, _id).Str
					arg1 := gjson.Get(record, args.Get("0").Str).Num
					arg2 := gjson.Get(record, args.Get("1").Str).Num
					val := arg1 / arg2

					mp[id] += val //.Num
					fieldCount[id]++

				}

			default:
				err = fmt.Errorf("unknown %s operator", op)
			}

			return false
		})

		if err != nil {
			return nil, err
		}
	}

	for fld, count := range fieldCount {
		mp[fld] = mp[fld] / count
	}

	return mp, nil
} // average

func count(field string, records []string) (mp map[string]int) {

	mp = map[string]int{}
	for _, record := range records {
		mp[gjson.Get(record, field).String()]++
	}

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

// reFields remove or rename fields in result data
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

func orderBy(param gjson.Result, reverse int, data []string) (list []string) {

	objects := []gjson.Result{}
	for _, v := range data {
		objects = append(objects, gjson.Parse(v))
	}
	// sort here

	// check type
	typ := 2 //objects[0].Get(param).Type
	fmt.Println("type is ", typ)

	if typ == 2 {
		list = sortNumber2(param, objects)
	}
	if typ == 3 {
		list = sortString("age", reverse, objects)
	}

	return list
}

type param struct {
	field string
	value int64
}

func sortNumber2(fields gjson.Result, list []gjson.Result) []string {
	var params []param

	fields.ForEach(func(key, val gjson.Result) bool {
		params = append(params, param{key.Str, val.Int()})
		return true
	})

	fmt.Println("Fields: ", fields) // {age:1, name:1}
	fmt.Println("params: ", params) // [{age,1}, {name, 1}]
	//lenListFields := len(listField) // 2

	max := len(list)
	var tmp gjson.Result

	element := list[0]

	fld := params[0].field // e.g age

	for max != 1 {
		for i := 1; i < max; i++ {
			if element.Get(fld).Num < list[i].Get(fld).Num {
				tmp = list[i]
				list[i] = element
				element = tmp
			}

			if element.Get(fld).Num == list[i].Get(fld).Num {
			} // fmt.Printf("\n %s, %s", element.Get("name"), element.Get("name"))
		}

		max--
		tmp = list[max]
		list[max] = element
		element = tmp

	}

	list[0] = element
	res := []string{}
	if params[0].value == 1 {
		fmt.Println("not reverse : ", fields.Get(params[0].field).Num)
		for i := 0; i < len(list); i++ {
			res = append(res, list[i].String())
		}
		return res
	}

	fmt.Println("reverse : ", params[0].value)
	for i := len(list) - 1; i >= 0; i-- {
		res = append(res, list[i].String())
	}

	return res
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
		}
		max--
		tmp = list[max]
		list[max] = element
		element = tmp

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

		}
		max--
		tmp = list[max]
		list[max] = element
		element = tmp

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

// end
