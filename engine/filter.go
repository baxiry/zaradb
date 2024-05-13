package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// gjson.Type :
// json-array:5, int:2, string:3

// match verifies that data matches the conditions
func match(filter, data string) (result bool, err error) {
	// TODO sould ber return syntax error if op unknown

	result = true

	gjson.Parse(filter).ForEach(func(qk, qv gjson.Result) bool {

		dv := gjson.Get(data, qk.String())
		//fmt.Println("q value type : ", qv.Type)

		if qk.String() == "$and" {
			// fmt.Println("query key is $and")
		}
		if qk.String() == "$or" {
			// fmt.Println("query key is $or")
		}

		if qv.Type == 5 { // 5:json, int:2, string:3
			qv.ForEach(func(sqk, sqv gjson.Result) bool {
				if sqv.Type == 3 { // 3:string,
					fmt.Println("here with: ", sqk.String())

					switch sqk.String() {
					case "$gt":
						if !(dv.String() > sqv.String()) {
							result = false
						}
						return result

					case "$lt":
						if !(dv.String() < sqv.String()) {
							result = false
						}
						return result

					case "$gte":
						if !(dv.String() >= sqv.String()) {
							result = false
						}
						return result

					case "$lte":
						if !(dv.String() <= sqv.String()) {
							result = false
						}
						return result

					case "$eq":
						if dv.String() != sqv.String() {
							result = false
						}
						return result
					case "$ne":
						if dv.String() == sqv.String() {
							result = false
						}
						return result

					default:
						err = fmt.Errorf("unknown %s operation", sqk.Value())
						//fmt.Println("..wher here", sqk.Value(), sqk.Type)
						result = false
						return result
					}

				}

				switch sqk.String() {
				case "$gt":
					if !(dv.Int() > sqv.Int()) {
						result = false
					}
					return result

				case "$lt":
					if !(dv.Int() < sqv.Int()) {
						result = false
					}
					return result

				case "$gte":
					if !(dv.Int() >= sqv.Int()) {
						result = false
					}
					return result

				case "$lte":
					if !(dv.Int() <= sqv.Int()) {
						result = false
					}
					return result

				case "$eq":
					if dv.Int() != sqv.Int() {
						result = false
					}
					return result

				case "$ne":
					if dv.Int() == sqv.Int() {
						result = false
					}
					return result
				case "$in": // in array
					//if sqv.IsArray() {}
					for _, v := range sqv.Array() {
						if dv.String() == v.String() {
							return result
						}
					}
					result = false
					return result

				case "$nin": // not in
					for _, v := range sqv.Array() {
						if dv.String() == v.String() {
							result = false
							return result
						}
					}
					return result

				default:
					fmt.Println("why iam here ?", sqk.String())
					err = fmt.Errorf("unknown %s operation", sqk.String())
					result = false
					return result
				}
			})

			match(qv.String(), dv.String())
			return result
		}

		if dv.String() != qv.String() {
			result = false
		}
		return result // if true keep iterating
	})
	return result, err
}
