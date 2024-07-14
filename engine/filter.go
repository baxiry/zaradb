package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func getIds(query gjson.Result) (string, error) {

	coll := query.Get("collection").Str
	if coll == "" {
		return "", fmt.Errorf("no collection")
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

	sub := query.Get("subQuery")
	if sub.Raw != "" {
		fmt.Println("sub.Row is : ", sub.Raw)
		ids, _ := getIds(sub)
		stmt += ` where rowid in (` + ids + `);`
		fmt.Println(stmt)
	}

	rows, err := db.db.Query(stmt)
	if err != nil {
		return "", fmt.Errorf("db.Query %s", err)
	}
	defer rows.Close()

	record := ""
	rowids := ""
	rowid := ""

	for rows.Next() {
		if limit == 0 {
			break
		}

		record = ""
		rowid = ""
		err := rows.Scan(&rowid, &record)
		if err != nil {
			return "", fmt.Errorf("row.Scan %s", err)
		}

		ok, err := match(mtch, record)
		if err != nil {
			return "", fmt.Errorf("match %s", err)
		}

		if ok {
			if skip != 0 {
				skip--
				continue
			}
			rowids += rowid + ","
			limit--
		}
	}

	if rowids == "" {
		return "", fmt.Errorf("zero value")
	}

	return rowids[:len(rowids)-1], nil
}

// gjson.Type :
// json:5, array:5, int:2, string:3

// match verifies that data matches the conditions
func match(filter gjson.Result, data string) (result bool, err error) {
	// TODO should return syntax error if op unknown

	result = true

	filter.ForEach(func(queryKey, queryVal gjson.Result) bool {

		dataVal := gjson.Get(data, queryKey.String())

		if queryVal.Type == 5 { // 5:json, int:2, string:3
			queryVal.ForEach(func(subQueryKey, subQueryVal gjson.Result) bool {

				if subQueryVal.Type == 3 { // 3:string,
					//fmt.Println("here with: ", subQueryKey.String())

					switch subQueryKey.Str {

					// comparition
					case "$gt":
						if !(dataVal.Str > subQueryVal.Str) {
							result = false
						}
						return result

					case "$lt":
						if !(dataVal.Str < subQueryVal.Str) {
							result = false
						}
						return result

					case "$gte":
						if !(dataVal.Str >= subQueryVal.Str) {
							result = false
						}
						return result

					case "$lte":
						if !(dataVal.Str <= subQueryVal.Str) {
							result = false
						}
						return result

					case "$eq":
						if dataVal.Str != subQueryVal.Str {
							result = false
						}
						return result
					case "$ne":
						if dataVal.Str == subQueryVal.Str {
							result = false
						}
						return result

					//
					case "$or": // not in

						result = false
						return result

					case "$st": // start with ..
						if !strings.HasPrefix(dataVal.Str, subQueryVal.Str) {
							result = false
						}
						return result

					case "$en": // end with ..
						if !strings.HasSuffix(dataVal.Str, subQueryVal.Str) {
							result = false
						}
						return result

					case "$c": // contains ..
						if !strings.Contains(dataVal.Str, subQueryVal.Str) {
							result = false
						}
						return result

					default:
						err = fmt.Errorf("unknown %s operation", subQueryKey.Value())
						//fmt.Println("..wher here", subQueryKey.Value(), subQueryKey.Type)
						result = false
						return result
					}
				}

				switch subQueryKey.Str {
				case "$gt":
					if !(dataVal.Num > subQueryVal.Num) {
						result = false
					}
					return result

				case "$lt":
					if !(dataVal.Num < subQueryVal.Num) {
						result = false
					}
					return result

				case "$gte":
					if !(dataVal.Num >= subQueryVal.Num) {
						result = false
					}
					return result

				case "$lte":
					if !(dataVal.Num <= subQueryVal.Num) {
						result = false
					}
					return result

				case "$eq":
					if dataVal.Num != subQueryVal.Num {
						result = false
					}
					return result

				case "$ne":
					if dataVal.Num == subQueryVal.Num {
						result = false
					}
					return result

				case "$in": // in array
					for _, v := range subQueryVal.Array() {
						//fmt.Println("subQueryVal", subQueryVal)
						if dataVal.Num == v.Num {
							return result
						}
					}
					result = false
					return result

				case "$nin": // not in
					for _, v := range subQueryVal.Array() {
						if dataVal.Num == v.Num {
							result = false
							return result
						}
					}
					return result

				default:

					// {$and:[{name:{$eq:"adam"}},{name:{$eq:"jawad"}}]}
					// {$or: [{name:{$eq:"adam"}}, {name:{$eq:"jawad"}}]}

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

					err = fmt.Errorf("unknown %s operation", subQueryKey.Str)
					result = false
					return result
				}
			})

			match(queryVal, queryVal.Str)
			return result
		}

		if dataVal.Str != queryVal.Str {
			result = false
		}
		return result // if true keep iterating
	})
	return result, err
}

// end
