package engine

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

// ....filter.....     ...data...    ...expected result...
func check(filter, data string, expected bool) string {
	filt := gjson.Parse(filter)
	result, _ := match(filt, data)

	if result != expected {
		return fmt.Sprintf(Yellow+"\nfilter:  %s\ndata:    %s\nexpected %v,\ngot:     %v\n"+Reset,
			filter, data, expected, result)
	}
	return ""
}

func Test_Match(t *testing.T) {

	// ....filter.....     ...data...    ...expected result...
	if res := check(`{"graded":true}`, `{"graded":true}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":true}`, `{"graded":false}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":false}`, `{"graded":false}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":false}`, `{"graded":true}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":null}`, `{"graded":null}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":null}`, `{"graded":true}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"graded":null}`, `{"graded":"null"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age": 18}`, `{"age": 18}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age": 18}`, `{"age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$eq":18}}`, `{"age": 18}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$eq":18}}`, `{"age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$ne": 18}}`, `{"age": 19}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$ne": 18}}`, `{"age": 18}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gt": 18}}`, `{"age": 19}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gt": 18}}`, `{"age": 18}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gt": 18}}`, `{"age": 17}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 18}}`, `{"age": 17}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 18}}`, `{"age": 18}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 18}}`, `{"age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gte": 18}}`, `{"age": 19}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gte": 18}}`, `{"age": 18}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$gte": 18}}`, `{"age": 17}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 18}}`, `{"age": 17}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 18}}`, `{"age": 18}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 18}}`, `{"age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 20}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 27}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 19}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lt": 28, "$gt": 18}}`, `{"age": 18}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 20}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 28}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 18}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$lte": 28, "$gte": 18}}`, `{"age": 16}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$in":[28,29,30]}}`, `{"age": 29}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$in":[28,29,30]}}`, `{"age": 9}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$nin":[28,29,30]}}`, `{"age": 9}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"age":{"$nin":[28,29,30]}}`, `{"age": 29}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"$and":[{"name":{"$eq":"adam"}}, {"age":{"$eq":29}} ] }`, `{"name":"adam","age": 29}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"$and":[{"name":{"$eq":"adam"}}, {"age":{"$eq":29}} ] }`, `{"name":"adam","age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"$or":[{"name":{"$eq":"adam"}}, {"age":{"$eq":29}} ] }`, `{"name":"jhon","age": 29}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"$or":[{"name":{"$eq":"adam"}}, {"age":{"$eq":29}} ] }`, `{"name":"jhon","age": 19}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":"adam"}`, `{"name":"adam"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":"adam"}`, `{"name":"kamal"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$eq":"adam"}}`, `{"name":"adam"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$eq":"adam"}}`, `{"name":"john"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$ne":"adam"}}`, `{"name":"john"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$ne":"adam"}}`, `{"name":"adam"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$st":"ad"}}`, `{"name":"adam"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$st":"ad"}}`, `{"name":"john"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$nst":"ad"}}`, `{"name":"john"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$nst":"ad"}}`, `{"name":"adam"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$en":"am"}}`, `{"name":"adam"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$en":"ad"}}`, `{"name":"john"}`, false); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$nen":"hn"}}`, `{"name":"adam"}`, true); res != "" {
		t.Error(res)
	}

	if res := check(`{"name":{"$nen":"am"}}`, `{"name":"adam"}`, false); res != "" {
		t.Error(res)
	}
	//
	if res := check(`{"name":{"$nen":"am"}}`, `{"name":"adam"}`, false); res != "" {
		t.Error(res)
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

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)
