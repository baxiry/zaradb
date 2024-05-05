package engine

import (
	"github.com/tidwall/gjson"
)

// match verifies that data matches the conditions
func match(filter, data string) (result bool) {
	// TODO sould ber return syntax error if op unknown

	result = true

	gjson.Parse(filter).ForEach(func(qk, qv gjson.Result) bool {

		dv := gjson.Get(data, qk.String())
		//fmt.Println("q value type : ", qv.Type)

		if qv.Type == 5 {
			//fmt.Println("inter ", qv.Type)
			qv.ForEach(func(sqk, sqv gjson.Result) bool {
				//	fmt.Println("    sqv type : ", sqv.Type)
				switch sqk.String() {
				case "$gt":
					if !(dv.Int() > sqv.Int()) {
						result = false
						return false
					}
					return result

				case "$lt":
					if !(dv.Int() < sqv.Int()) {
						result = false
						return false
					}
					return result

				case "$gte":
					if !(dv.Int() >= sqv.Int()) {
						result = false
						return false
					}
					return result

				case "$lte":
					if !(dv.Int() <= sqv.Int()) {
						result = false
						return false
					}
					return result

				case "$eq":
					if dv.Int() != sqv.Int() {
						result = false
						return false
					}
					return result
				default:
					result = false
					return result
				}
			})

			match(qv.String(), dv.String())
			return result
		}

		if dv.String() != qv.String() {
			result = false
			return result
		}
		return result // if true keep iterating
	})
	return result
}
