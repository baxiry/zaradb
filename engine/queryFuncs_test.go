package engine

import (
	"fmt"
	"os"
	"testing"

	"github.com/tidwall/gjson"
)

var got string

func Test_NewDB(t *testing.T) {
	s = NewDB("tmptest.db")
	if s == nil {
		t.Error("Store should not be nil")
	}
}

func Test_insertOne(t *testing.T) {
	json := `{"collection":"test", "action":"insertOne","data":{"name":"adam", "age": 23}}`
	query := gjson.Parse(json)
	s.insertOne(query)

	res := s.db.QueryRow("select obj from test where rowid=1")
	res.Scan(&got)

	exp := `{"_id":1, "name":"adam", "age": 23}`
	if got != exp {
		t.Errorf("\n%sexpect %s\ngot %s %s", Yellow, exp, got, Reset)
	}

}

func Test_findOne(t *testing.T) {
	query := `{"collection":"test", "action":"findOne","match":{"name":"adam"}}` //

	res := HandleQueries(query)
	exp := `{"_id":1, "name":"adam", "age": 23}`

	if res != exp {
		t.Errorf("expect %s\ngot %s", exp, res)
	}
}

func Test_HandleQury(t *testing.T) {
	// TODO
	_ = HandleQueries("")
}

func cfindById(jsonQuery, expected string) string {
	query := gjson.Parse(jsonQuery)

	result := s.findById(query)

	if result != expected {
		return fmt.Sprintf(Yellow+"\nexp\t %s\ngot\t %s\n"+Reset, expected, result)
	}

	return ""
}

func Test_findById(t *testing.T) {

	// "Valid collection and ID"
	if res := cfindById(`{"collection":"test", "action":"findById", "_id":1}`, `{"_id":1, "name":"adam", "age": 23}`); res != "" {
		t.Error("\nValid collection and ID", res)
	}

	// "Collection does not exist"
	if res := cfindById(`{"collection":"unknown", "action":"findById", "_id":1}`, `{"error": "collection does not exist"}`); res != "" {
		t.Error("\nCollection does not exist", res)
	}

	// "ID does not exist"
	// "Collection does not exist"
	if res := cfindById(`{"collection":"test", "action":"findById", "_id":123}`, `{"error": "_id does not exist"}`); res != "" {
		t.Error("\nID does not exist", res)
	}
}

func Test_insertMany(t *testing.T) {
	json := `{"collection":"testInsertMany", "action":"insertMany","data":[{"name":"adam1", "age": 21},{"name":"adam2", "age": 22},{"name":"adam3", "age": 23}]}`

	if s.insertMany(gjson.Parse(json)) != `{"ak":"insertMany Done"}` {
		t.Error("error when insertMany")
	}

	exp := `{"_id":1,"name":"adam1", "age": 21}`
	res := s.db.QueryRow("select obj from testInsertMany where rowid=1")
	res.Scan(&got)
	if got != exp {
		t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
	}

	exp = `{"_id":2,"name":"adam2", "age": 22}`
	res = s.db.QueryRow("select obj from testInsertMany where rowid=2")
	res.Scan(&got)
	if got != exp {
		t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
	}

	exp = `{"_id":3,"name":"adam3", "age": 23}`
	res = s.db.QueryRow("select obj from testInsertMany where rowid=3")
	res.Scan(&got)

	if got != exp {
		t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
	}
}

func Test_findMany(t *testing.T) {

	got := s.findMany(gjson.Parse(`{"collection":"testInsertMany", "action":"findMany"}`))
	exp := `[{"_id":1,"name":"adam1", "age": 21},{"_id":2,"name":"adam2", "age": 22},{"_id":3,"name":"adam3", "age": 23}]`
	if got != exp {
		t.Errorf(Red+"got %s\nexp %s"+Reset, got, exp)
	}

	got = s.findMany(gjson.Parse(`{"collection":"testInsertMany", "action":"findMany", "limit": 1}`))
	exp = `[{"_id":1,"name":"adam1", "age": 21}]`
	if got != exp {
		t.Errorf(Red+"got %s\nexp %s"+Reset, got, exp)
	}

	got = s.findMany(gjson.Parse(`{"collection":"testInsertMany", "action":"findMany", "limit": 1, "skip", 2}`))
	exp = `[{"_id":3,"name":"adam3", "age": 23}]`
	if got != exp {
		t.Errorf(Red+"got %s\nexp %s"+Reset, got, exp)
	}

	got = s.findMany(gjson.Parse(`{"collection":"testInsertMany", "action":"findMany", "skip", 3}`))
	exp = `[]`
	if got != exp {
		t.Errorf(Red+"got %s\nexp %s"+Reset, got, exp)
	}

}

func Test_updateOne(t *testing.T) {

	ak := s.updateOne(gjson.Parse(`{"collection":"test", "action":"updateOne", "data":{"age":26}}, "match":{"name":"adam"}`))
	if ak != `{"ak": "update: done"}` {
		fmt.Println("akk: ", ak)
		t.Errorf(Red+"got\t %s\nexp\t %s"+Reset, ak, `{"ak": "update: done"}`)
	}

	exp := `{"_id":1,"name":"adam","age":26}`
	res := s.db.QueryRow("select obj from test;")
	res.Scan(&got)

	if got != exp {
		t.Errorf("\n%sexp\t %s\ngot\t %s %s", Yellow, exp, got, Reset)
	}
}

func Test_Close(t *testing.T) {
	err := s.db.Close()

	if err != nil {
		t.Error("Store should be nil")
	}
	err = s.db.Ping()
	if err.Error() != "sql: database is closed" {
		t.Error("store steal work")
	}
	os.Remove("tmptest.db")
}

/*

	testCases := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "Collection does not exist",
			query:    `{"collection":"unknown","_id":1}`,
			expected: `{"error": "collection does not exist"}`,
		},
		{
			name:     "ID does not exist",
			query:    `{"collection":"test","_id":123}`,
			expected: `{"error": "_id does not exist"}`, // No value in DB for this key
		},
		{
			name:     "Valid collection and ID",
			query:    `{"collection":"test","_id":1}`,
			expected: `{"_id":1, "name":"adam", "age": 23}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := gjson.Parse(tc.query)
			result := s.findById(query)

			if result != tc.expected {
				t.Errorf(Yellow+"\nexpected %s\n got %s\n"+Reset, tc.expected, result)
			}
		})
	}

*/
