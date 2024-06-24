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

					switch subQueryKey.String() {
					case "$gt":
						if !(dataVal.String() > subQueryVal.String()) {
							result = false
						}
						return result

					case "$lt":
						if !(dataVal.String() < subQueryVal.String()) {
							result = false
						}
						return result

					case "$gte":
						if !(dataVal.String() >= subQueryVal.String()) {
							result = false
						}
						return result

					case "$lte":
						if !(dataVal.String() <= subQueryVal.String()) {
							result = false
						}
						return result

					case "$eq":
						if dataVal.String() != subQueryVal.String() {
							result = false
						}
						return result
					case "$ne":
						if dataVal.String() == subQueryVal.String() {
							result = false
						}
						return result

					case "$or": // not in

						result = false
						return result

					case "$st": // is start with
						if !strings.HasPrefix(dataVal.String(), subQueryVal.String()) {
							result = false
						}
						return result

					case "$en": // is end with
						if !strings.HasSuffix(dataVal.String(), subQueryVal.String()) {
							result = false
						}
						return result

					case "$c": // is contains
						if !strings.Contains(dataVal.String(), subQueryVal.String()) {
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

				switch subQueryKey.String() {
				case "$gt":
					if !(dataVal.Int() > subQueryVal.Int()) {
						result = false
					}
					return result

				case "$lt":
					if !(dataVal.Int() < subQueryVal.Int()) {
						result = false
					}
					return result

				case "$gte":
					if !(dataVal.Int() >= subQueryVal.Int()) {
						result = false
					}
					return result

				case "$lte":
					if !(dataVal.Int() <= subQueryVal.Int()) {
						result = false
					}
					return result

				case "$eq":
					if dataVal.Int() != subQueryVal.Int() {
						result = false
					}
					return result

				case "$ne":
					if dataVal.Int() == subQueryVal.Int() {
						result = false
					}
					return result
				case "$in": // in array
					for _, v := range subQueryVal.Array() {
						if dataVal.String() == v.String() {
							return result
						}
					}
					result = false
					return result

				case "$nin": // not in
					for _, v := range subQueryVal.Array() {
						if dataVal.String() == v.String() {
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

					err = fmt.Errorf("unknown %s operation", subQueryKey.String())
					result = false
					return result
				}
			})

			match(queryVal, queryVal.String())
			return result
		}

		if dataVal.String() != queryVal.String() {
			result = false
		}
		return result // if true keep iterating
	})
	return result, err
}
