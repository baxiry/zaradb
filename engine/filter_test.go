package engine

import (
	"testing"

	"github.com/tidwall/gjson"
)

var testCases = []struct {
	filter   string
	data     string
	expected bool
}{
	{`{"name":"adam"}`, `{"name":"adam"}`, true},
	{`{"name":"adam"}`, `{"name":"kamal"}`, false},
}

func TestMatch(t *testing.T) {
	// الحالات الاختبارية

	// تنفيذ الحالات الاختبارية باستخدام t.Run
	for _, tcase := range testCases {
		result, _ := match(gjson.Parse(tcase.filter), tcase.data)
		if result != tcase.expected {
			t.Errorf("For input %q and pattern %q, expected %v but got %v", tcase.filter, tcase.data, tcase.expected, result)
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
