package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// gjson.Type =>  json:5, array:5, int:2, string:3

// match verifies that data matches the conditions
func match(filter gjson.Result, data string, ids ...int64) (result bool, err error) {
	// TODO should return syntax error if op unknown

	result = true

	filter.ForEach(func(queryKey, queryVal gjson.Result) bool {

		dataVal := gjson.Get(data, queryKey.Str)

		if queryVal.Type == 5 { // 5:json
			// {name:{$eq:"adam"}, age:{$gt: 18}}

			queryVal.ForEach(func(sQueryKey, sQueryVal gjson.Result) bool {

				if sQueryVal.Type == 3 { // 3:string,
					// from :  {$eq:"adam"} , sQueryKey is $eq, sQueryVal is "adam"

					switch sQueryKey.Str {

					// compare sQueryKey
					case "$gt":
						if !(dataVal.Str > sQueryVal.Str) {
							result = false
						}
						return result

					case "$lt":
						if !(dataVal.Str < sQueryVal.Str) {
							result = false
						}
						return result

					case "$gte":
						if !(dataVal.Str >= sQueryVal.Str) {
							result = false
						}
						return result

					case "$lte":
						if !(dataVal.Str <= sQueryVal.Str) {
							result = false
						}
						return result

					case "$eq":
						if dataVal.Str != sQueryVal.Str {
							result = false
						}
						return result
					case "$ne":
						if dataVal.Str == sQueryVal.Str {
							result = false
						}
						return result

					case "$st": // start with ..
						if !strings.HasPrefix(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					case "$en": // end with ..
						if !strings.HasSuffix(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					case "$c": // contains ..
						if !strings.Contains(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					case "$nst": // start with ..
						if strings.HasPrefix(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					case "$nen": // end with ..
						if strings.HasSuffix(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					case "$nc": // contains ..
						if strings.Contains(dataVal.Str, sQueryVal.Str) {
							result = false
						}
						return result

					default:

						err = fmt.Errorf("unknown %s operation", sQueryKey.Str)
						//fmt.Println("..wher here", sQueryKey.Value(), sQueryKey.Type)
						result = false
						return result
					}
				}

				// if sQueryVal is number
				switch sQueryKey.Str {

				case "$gt":
					if !(dataVal.Num > sQueryVal.Num) {
						result = false
					}
					return result

				case "$lt":
					if !(dataVal.Num < sQueryVal.Num) {
						result = false
					}
					return result

				case "$gte":
					if !(dataVal.Num >= sQueryVal.Num) {
						result = false
					}
					return result

				case "$lte":
					if !(dataVal.Num <= sQueryVal.Num) {
						result = false
					}
					return result

				case "$eq":
					if dataVal.Num != sQueryVal.Num {
						result = false
					}
					return result

				case "$ne":
					if dataVal.Num == sQueryVal.Num {
						result = false
					}
					return result

				case "$in": // in array

					// handle Str arr
					if sQueryKey.Type == 3 {
						for _, v := range sQueryVal.Array() {
							//fmt.Println("sQueryVal", sQueryVal)
							if dataVal.Str == v.Str {
								return result
							}
						}
						result = false
						return result
					}

					// handle Num arr
					for _, v := range sQueryVal.Array() {
						//fmt.Println("sQueryVal", sQueryVal)
						if dataVal.Num == v.Num {
							return result
						}
					}
					result = false
					return result

				case "$nin": // not in
					// handle string arr
					if sQueryKey.Type == 3 {
						for _, v := range sQueryVal.Array() {
							if dataVal.Str == v.Str {
								result = false
								return result
							}
						}
						return result
					}

					// handle Num arr
					for _, v := range sQueryVal.Array() {
						if dataVal.Num == v.Num {
							result = false
							return result
						}
					}
					return result

				default:

					// {$and:[{name:{$eq:"adam"}},{name:{$eq:"jawad"}}]}
					if queryKey.Str == "$and" {

						for _, v := range queryVal.Array() {
							res, _ := match(v, data)
							if !res {
								result = false
								return result
							}
						}
						return result
					}

					// {$or: [{name:{$eq:"adam"}}, {name:{$eq:"jawad"}}]}
					if queryKey.Str == "$or" {

						for _, v := range queryVal.Array() {

							res, _ := match(v, data)
							if res {
								return result
							}
						}
						result = false
						return result
					}

					err = fmt.Errorf("unknown %s operation", sQueryKey.Str)
					result = false
					return result
				}
			})

			match(queryVal, queryVal.Str)
			return result
		}

		// if queryVal is number : {age: 10}
		if queryVal.Type == 2 {
			if queryVal.Num != dataVal.Num {
				result = false
			}
		}

		// if queryVal is string : {name: "adam"}
		if queryVal.Str != dataVal.Str {
			result = false
		}

		// if result is true then keep iterating
		return result
	})
	return result, err
}

func getIds(query gjson.Result) (ids []int64) {

	coll := query.Get("collection").Str
	if coll == "" {
		return nil
	}

	mtch := query.Get("match")

	if mtch.String() == "" {
		fmt.Println("match.Str is empty")
	}

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		limit = 100
	}

	stmt := `select rowid, record from ` + coll

	sub := query.Get("sQuery")

	if sub.Raw != "" {
		fmt.Println("sub.Row is : ", sub.Raw)
		ids = getIds(sub)
	}

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil
	}
	defer rows.Close()

	record := ""
	rowid := 0

	for rows.Next() {
		if limit == 0 {
			break
		}

		record = ""
		rowid = 0
		_ = rows.Scan(&rowid, &record)

		ok, err := match(mtch, record)
		if err != nil {
			fmt.Printf("match %s\n", err)
			return nil
		}

		if ok {
			if skip != 0 {
				skip--
				continue
			}
			ids = append(ids, int64(rowid))
			limit--
		}
	}
	fmt.Println("\n", ids)

	return ids
}

// end
