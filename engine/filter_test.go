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
	{`{"match":{"age": 18}}`, `{"data":{"age": 18}}`, true},
	{`{"match":{"age": 18}}`, `{"data":{"age": 19}}`, false},

	{`{"match":{"age":{"$eq":18}}}`, `{"data":{"age": 18}}`, true},
	{`{"match":{"age":{"$eq":18}}}`, `{"data":{"age": 19}}`, false},

	{`{"match":{"age":{"$ne": 18}}}`, `{"data":{"age": 19}}`, true},
	{`{"match":{"age":{"$ne": 18}}}`, `{"data":{"age": 18}}`, false},

	{`{"match":{"age":{"$gt": 18}}}`, `{"data":{"age": 19}}`, true},
	{`{"match":{"age":{"$gt": 18}}}`, `{"data":{"age": 18}}`, false},
	{`{"match":{"age":{"$gt": 18}}}`, `{"data":{"age": 17}}`, false},

	{`{"match":{"age":{"$lt": 18}}}`, `{"data":{"age": 17}}`, true},
	{`{"match":{"age":{"$lt": 18}}}`, `{"data":{"age": 18}}`, false},
	{`{"match":{"age":{"$lt": 18}}}`, `{"data":{"age": 19}}`, false},

	{`{"match":{"age":{"$gte": 18}}}`, `{"data":{"age": 19}}`, true},
	{`{"match":{"age":{"$gte": 18}}}`, `{"data":{"age": 18}}`, true},
	{`{"match":{"age":{"$gte": 18}}}`, `{"data":{"age": 17}}`, false},

	{`{"match":{"age":{"$lte": 18}}}`, `{"data":{"age": 17}}`, true},
	{`{"match":{"age":{"$lte": 18}}}`, `{"data":{"age": 18}}`, true},
	{`{"match":{"age":{"$lte": 18}}}`, `{"data":{"age": 19}}`, false},

	{`{"match":{"age":{"$lt": 28, "$gt": 18}}}`, `{"data":{"age": 20}}`, true},
	{`{"match":{"age":{"$lt": 28, "$gt": 18}}}`, `{"data":{"age": 27}}`, true},
	{`{"match":{"age":{"$lt": 28, "$gt": 18}}}`, `{"data":{"age": 19}}`, true},
	{`{"match":{"age":{"$lt": 28, "$gt": 18}}}`, `{"data":{"age": 18}}`, false},

	{`{"match":{"age":{"$lte": 28, "$gte": 18}}}`, `{"data":{"age": 20}}`, true},
	{`{"match":{"age":{"$lte": 28, "$gte": 18}}}`, `{"data":{"age": 28}}`, true},
	{`{"match":{"age":{"$lte": 28, "$gte": 18}}}`, `{"data":{"age": 18}}`, true},
	{`{"match":{"age":{"$lte": 28, "$gte": 18}}}`, `{"data":{"age": 16}}`, false},

	{`{"match":{"age":{"$in":[28,29,30]}}}`, `{"data":{"age": 29}}`, true},
	{`{"match":{"age":{"$in":[28,29,30]}}}`, `{"data":{"age": 9}}`, false},

	{`{"match":{"age":{"$nin":[28,29,30]}}}`, `{"data":{"age": 9}}`, true},
	{`{"match":{"age":{"$nin":[28,29,30]}}}`, `{"data":{"age": 29}}`, false},
	// string cases
	{`{"match":{"name":"adam"}}`, `{"data":{"name":"adam"}}`, true},
	{`{"match":{"name":"adam"}}`, `{"data":{"name":"kamal"}}`, false},

	{`{"match":{"name":{"$eq":"adam"}}}`, `{"data":{"name":"adam"}}`, true},
	{`{"match":{"name":{"$eq":"adam"}}}`, `{"data":{"name":"john"}}`, false},

	{`{"match":{"name":{"$ne":"adam"}}}`, `{"data":{"name":"john"}`, true},
	{`{"match":{"name":{"$ne":"adam"}}}`, `{"data":{"name":"adam"}}`, false},
}

func Test_Match(t *testing.T) {

	for _, tcase := range testCases {
		filt := gjson.Get(tcase.filter, "match")
		data := gjson.Get(tcase.data, "data").Raw
		result, _ := match(filt, data)

		if result != tcase.isMatch {
			t.Errorf("\nfilter:  %s\ndata:    %s\nexpected %v,\ngot:     %v\n",
				filt, data, tcase.isMatch, result)
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
