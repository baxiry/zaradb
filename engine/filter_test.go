package engine

import (
	"testing"

	"github.com/tidwall/gjson"
)

var testCases = []struct {
	filter  string
	data    string
	isMatch bool
}{
	// ....filter.....     ...data...    ...isMatch...
	//  cases of numbers
	{`{"age": 18}`, `{"age": 18}`, true},
	{`{"age": 18}`, `{"age": 19}`, false},

	{`{"age":{"$eq":18}}`, `{"age": 18}`, true},
	{`{"age":{"$eq":18}}`, `{"age": 19}`, false},

	{`{"age":{"$ne": 18}}`, `{"age": 19}`, true},
	{`{"age":{"$ne": 18}}`, `{"age": 18}`, false},

	{`{"age":{"$gt": 18}}`, `{"age": 19}`, true},
	{`{"age":{"$gt": 18}}`, `{"age": 18}`, false},
	{`{"age":{"$gt": 18}}`, `{"age": 17}`, false},

	{`{"age":{"$lt": 18}}`, `{"age": 17}`, true},
	{`{"age":{"$lt": 18}}`, `{"age": 18}`, false},
	{`{"age":{"$lt": 18}}`, `{"age": 19}`, false},

	{`{"age":{"$gte": 18}}`, `{"age": 19}`, true},
	{`{"age":{"$gte": 18}}`, `{"age": 18}`, true},
	{`{"age":{"$gte": 18}}`, `{"age": 17}`, false},

	{`{"age":{"$lte": 18}}`, `{"age": 17}`, true},
	{`{"age":{"$lte": 18}}`, `{"age": 18}`, true},
	{`{"age":{"$lte": 18}}`, `{"age": 19}`, false},

	{`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 20}`, true},
	{`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 27}`, true},
	{`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 19}`, true},
	{`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 18}`, false},

	{`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 20}`, true},
	{`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 28}`, true},
	{`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 18}`, true},
	{`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 16}`, false},

	{`{"age":{"$in":[28,29,30]}}`, `{"age": 29}`, true},
	{`{"age":{"$in":[28,29,30]}}`, `{"age": 9}`, false},

	{`{"age":{"$nin":[28,29,30]}}`, `{"age": 9}`, true},
	{`{"age":{"$nin":[28,29,30]}}`, `{"age": 29}`, false},
	// string cases
	{`{"name":"adam"}`, `{"name":"adam"}`, true},
	{`{"name":"adam"}`, `{"name":"kamal"}`, false},

	{`{"name":{"$eq":"adam"}}`, `{"name":"adam"}`, true},
	{`{"name":{"$eq":"adam"}}`, `{"name":"john"}`, false},

	{`{"name":{"$ne":"adam"}}`, `{"name":"john"}`, true},
	{`{"name":{"$ne":"adam"}}`, `{"name":"adam"}`, false},
}

func Test_Match(t *testing.T) {

	for _, tcase := range testCases {
		result, _ := match(gjson.Parse(tcase.filter), tcase.data)
		if result != tcase.isMatch {
			t.Errorf("\n\nfilter:  %s\ndata:    %s\nexpected %v,\ngot:     %v\n",
				tcase.filter, tcase.data, tcase.isMatch, result)
		}
	}
}

var data = []string{

	`{"_id": 2, "name":"karim", "age":37, "contact":{"email": "karim@email.com", "tele": "00234856781"}}`,
	`{"_id": 3, "name":"jawad", "age":24, "contact":{"email": "jawad@email.com", "tele": "00123845672"}}`,
	`{"_id": 4, "name":"imane", "age":26, "contact":{"email": "imane@email.com", "tele": "00239850303"}}`,
	`{"_id": 5, "name":"adams", "age":31, "contact":{"email": "adams@email.com", "tele": "00234850004"}}`,
	`{"_id": 6, "name":"johns", "age":11, "contact":{"email": "adams@email.com", "tele": "00234850105"}}`,
	`{"_id": 7, "name":"jawar", "age":10, "contact":{"email": "jawar@email.com", "tele": "00234850226"}}`,
	`{"_id": 8, "name":"kadir", "age":13, "contact":{"email": "kadir@email.com", "tele": "00238450337"}}`,
	`{"_id": 9, "name":"hamid", "age":16, "contact":{"email": "hamid@email.com", "tele": "00238450048"}}`,
	`{"_id":10, "name":"rajab", "age":18, "contact":{"email": "rajab@email.com", "tele": "00230450059"}}`,
	`{"_id":11, "name":"samir", "age":19, "contact":{"email": "samir@email.com", "tele": "00230450060"}}`,
	`{"_id":12, "name":"akram", "age":21, "contact":{"email": "akram@email.com", "tele": "00239450071"}}`,
	`{"_id":13, "name":"rabih", "age":31, "contact":{"email": "rabih@email.com", "tele": "00230450083"}}`,
	`{"_id":14, "name":"monir", "age":31, "contact":{"email": "monir@email.com", "tele": "00239450094"}}`,
}
