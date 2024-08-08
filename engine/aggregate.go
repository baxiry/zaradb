package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const siparator = "_:_"

// not implemented yet
func min(query gjson.Result) (int, error)     { return 0, nil }
func max(query gjson.Result) (int, error)     { return 0, nil }
func sum(query gjson.Result) (int, error)     { return 0, nil }
func average(query gjson.Result) (int, error) { return 0, nil }
func count(query gjson.Result) (int, error) {

	// TODO parse hol qury one time
	coll := query.Get("collection").Str
	if coll == "" {
		return 0, errors.New("forgot collection name")
	}

	mtch := query.Get("match")

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		limit = 100 // what is default setting ?
	}

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		fmt.Println("err", err.Error())
		return 0, err
	}
	defer rows.Close()

	record := ""
	total := 0

	for rows.Next() {

		if limit == 0 {
			break
		}
		if skip != 0 {
			skip--
			continue
		}

		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return 0, err
		}

		ok, err := match(mtch, record)
		if err != nil {
			return 0, err
		}

		if ok {
			total++
			limit--
		}
	}

	return total, nil
}

// not implemented yet
func aggrigate(query gjson.Result) string {

	// TODO parse hol qury one time
	coll := query.Get("collection").Str
	if coll == "" {
		return `{"error":"forgot collection name "}`
	}

	mtch := query.Get("match")

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		limit = 100 // what is default setting ?
	}

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	record := ""
	listData := make([]string, 0)

	for rows.Next() {

		if limit == 0 {
			break
		}
		if skip != 0 {
			skip--
			continue
		}

		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return err.Error() // TODO standaring errors
		}

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}

		if ok {
			listData = append(listData, record)
			limit--
		}
	}

	// order :
	order := query.Get("orderBy").Str
	reverse := query.Get("reverse").Int()

	if order != "" {
		listData = orderBy(order, int(reverse), listData)
	}

	// TODO aggrigate here

	// remove or rename some fields
	flds := query.Get("fields")
	listData = reFields(listData, flds)

	records := "["

	for i := 0; i < len(listData); i++ {
		records += listData[i] + ","
	}

	if len(records) == 1 {
		return records + "]"
	}

	return records[:len(records)-1] + "]"

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
