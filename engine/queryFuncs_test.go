package engine

import (
	"os"
	"testing"

	"github.com/tidwall/gjson"
	"go.etcd.io/bbolt"
)

var s *Store

//func Test_NewDB(){}
//func Test_findOne(){}
//func Test_findById(){}

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

	s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("test"))
		res := bucket.Get(uint64ToBytes(1))
		exp := `{"_id":1, "name":"adam", "age": 23}`
		if string(res) != exp {
			t.Errorf("\n%sexpect %s\ngot %s %s", Yellow, exp, res, Reset)
		}

		return nil
	})
}

func Test_findOne(t *testing.T) {
	json := `{"collection":"test", "action":"findOne","match":{"name":"adam"}}` //
	query := gjson.Parse(json)

	res := s.findOne(query)
	exp := []byte(`{"_id":1, "name":"adam", "age": 23}`)

	if string(res) != string(exp) {
		t.Errorf("expect %s\ngot %s", exp, res)
	}
}

func Test_findById(t *testing.T) {

	testCases := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "Valid collection and ID",
			query:    `{"collection":"test","_id":1}`,
			expected: `{"_id":1, "name":"adam", "age": 23}`,
		},
		{
			name:     "Collection does not exist",
			query:    `{"collection":"unknown","_id":1}`,
			expected: `{"error": "collection unknown not exist"}`,
		},
		{
			name:     "ID does not exist",
			query:    `{"collection":"test","_id":"nonexistent"}`,
			expected: "", //`{"error": ""}`, // No value in DB for this key
		},
		{
			name:     "ID does not exist",
			query:    `{"collection":"test","_id":123}`,
			expected: "", //`{"error": ""}`, // No value in DB for this key
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := gjson.Parse(tc.query)
			result := s.findById(query)

			if result != tc.expected {
				t.Errorf(Yellow+"expected %s, got %s"+Reset, tc.expected, result)
			}
		})
	}
}

func Test_insertMany(t *testing.T) {
	json := `{"collection":"insertTest", "action":"insertMany","data":[{"name":"adam1", "age": 21},{"name":"adam2", "age": 22},{"name":"adam3", "age": 23}]}`
	query := gjson.Parse(json)
	s.insertMany(query)

	s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("insertTest"))

		exp := `{"_id":1,"name":"adam1", "age": 21}`
		got := string(bucket.Get(uint64ToBytes(1)))
		if got != exp {
			t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
		}

		exp = `{"_id":2,"name":"adam2", "age": 22}`

		got = string(bucket.Get(uint64ToBytes(2)))
		if got != exp {
			t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
		}

		exp = `{"_id":3,"name":"adam3", "age": 23}`

		got = string(bucket.Get(uint64ToBytes(3)))
		if got != exp {
			t.Errorf("\n%sgot %s\nexp %s %s", Yellow, got, exp, Reset)
		}

		return nil
	})
}

func Test_findMany(t *testing.T) {
	json := `{"collection":"insertTest", "action":"findMany"}` //,"match":{"name":"adam"}
	query := gjson.Parse(json)

	got := s.findMany(query)
	exp := `[{"_id":1,"name":"adam1", "age": 21},{"_id":2,"name":"adam2", "age": 22},{"_id":3,"name":"adam3", "age": 23}]`

	if got != exp {
		t.Errorf("got %s\nexp %s", got, exp)
	}

}

func Test_Close(t *testing.T) {
	err := s.db.Close()

	if err != nil {
		t.Error("Store should be nil")
	}
	os.Remove("tmptest.db")
}
