package engine

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

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

					switch subQueryKey.Str { // .String()

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

					case "$st": // is it start with ?
						if !strings.HasPrefix(dataVal.Str, subQueryVal.Str) {
							result = false
						}
						return result

					case "$en": // is it end with
						if !strings.HasSuffix(dataVal.Str, subQueryVal.Str) {
							result = false
						}
						return result

					case "$c": // is it contains
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
						if dataVal.Str == v.Str {
							return result
						}
					}
					result = false
					return result

				case "$nin": // not in
					for _, v := range subQueryVal.Array() {
						if dataVal.Str == v.Str {
							result = false
							return result
						}
					}
					return result

				//case "$ins": // is it is sub query ?

				default:
					if subQueryKey.Str == "$ins" {
						fmt.Println("queryVal : ", subQueryVal.Raw)
						fmt.Println("ids: ", getIds(subQueryVal))
					}

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
