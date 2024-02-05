package store

import (
	"fmt"
	"os"
	"testing"
)

var db *Database
var coll *Collection

func Test_NewDatabase(t *testing.T) {
	var path = "test/test"

	db = NewDatabase("test")

	coll, _ = db.Collections["test"]

	_, err := os.Stat(path + "/00000000000000000001")
	if err != nil {
		if !os.IsExist(err) {
			t.Errorf("%s  should be exists", path)
		}
	}

}

func Test_get(t *testing.T) {
	var i uint64

	for i = 1; i < 12; i++ {
		coll.insert("hello")
	}

	l, _ := coll.log.LastIndex()
	for i = 1; i < l; i++ {
		data, err := coll.get(i)
		if err != nil {
			panic(err)
		}
		if data != fmt.Sprint(i)+" hello" {
			t.Error("data shoul be:", fmt.Sprint(i)+" hello")
		}
		t.Log(data, ".. ok")
	}

	data, err := coll.get(1)
	if err != nil {
		t.Log(err)
		t.Log(coll.lastIndex)
	}
	if data != "1 hello" {
		t.Error(data, "should be '1 hello'")
		t.Log(data)
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
