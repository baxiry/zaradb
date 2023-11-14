package dblite

import (
	"github.com/tidwall/gjson"
)

// valid verifies that data matches the conditions
func valid(query, data string) (result bool) {
	result = true

	gjson.Parse(query).ForEach(func(qk, qv gjson.Result) bool {

		dv := gjson.Get(data, qk.String())

		if qv.Type == 5 {
			qv.ForEach(func(sqk, sqv gjson.Result) bool {
				if sqk.String()[0] == '$' {
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
						// ??
					}
				}
				return result
			})

			valid(qv.String(), dv.String())
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
