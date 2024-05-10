package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// match verifies that data matches the conditions
func match(filter, data string) (result bool, err error) {
	// TODO sould ber return syntax error if op unknown

	result = true

	gjson.Parse(filter).ForEach(func(qk, qv gjson.Result) bool {

		dv := gjson.Get(data, qk.String())
		//fmt.Println("q value type : ", qv.Type)

		if qv.Type == 5 {
			//fmt.Println("inter ", qv.Type)
			qv.ForEach(func(sqk, sqv gjson.Result) bool {

				if sqv.Type == 3 {

					switch sqk.String() {
					case "$gt":
						if !(dv.String() > sqv.String()) {
							result = false
							return false
						}
						return result

					case "$lt":
						if !(dv.String() < sqv.String()) {
							result = false
							return false
						}
						return result

					case "$gte":
						if !(dv.String() >= sqv.String()) {
							result = false
							return false
						}
						return result

					case "$lte":
						if !(dv.String() <= sqv.String()) {
							result = false
							return false
						}
						return result

					case "$eq":
						if dv.String() != sqv.String() {
							result = false
							return false
						}
						return result
					case "$ne":
						if dv.String() == sqv.String() {
							result = false
							return false
						}
						return result

					default:
						err = fmt.Errorf("unknown %s opiration", sqk.String())
						result = false
						return result
					}

				}

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

				case "$ne":
					if dv.Int() == sqv.Int() {
						result = false
						return false
					}
					return result

				default:

					err = fmt.Errorf("unknown %s opiration", sqk.String())
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
	return result, err
}
