package store

import (
	"os"
	"testing"
)

var db *Database
var coll *Collection

func Test_NewDatabase(t *testing.T) {
	var path = "test"

	db = NewDatabase(path)

	coll, _ = db.Collections["test"]

	_, err := os.Stat(db.path + "/test" + "/00000000000000000001")
	if err != nil {
		if !os.IsExist(err) {
			t.Errorf("%s  should be exists", db.path+"/00000000000000000001")
		}

	}

}

func Test_insert(t *testing.T) {
	var i uint64

	data := "hello"

	for i = 1; i < 12; i++ {
		err := coll.insert(data)
		if err != nil {
			panic(err)
		}
	}

	l, _ := coll.log.LastIndex()
	for i = 1; i < l; i++ {
		data, err := coll.get(i)
		if err != nil {
			panic(err)
		}

		if data != "hello" {
			t.Errorf("shoul be: %s, not %s", "hello", data)
		}
	}

}

func Test_get(t *testing.T) {
	var i uint64

	want := "hello"

	l, _ := coll.log.LastIndex()
	for i = 1; i < l; i++ {
		data, err := coll.get(i)
		if err != nil {
			panic(err)
		}

		if data != want {
			t.Errorf("shoul be: %s, not %s", want, data)
		}
	}

}

func Test_Close(t *testing.T) {
	if len(db.Collections) == 0 {
		t.Error("db should ont be empty")
	}

	db.Close()

	if len(db.Collections) != 0 {
		t.Error("db should be empty")
	}

	if _, ok := db.Collections["test"]; ok {
		t.Error("db is not closed")
	}
}
