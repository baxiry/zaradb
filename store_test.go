package store

import (
	"fmt"
	"os"
	"strings"
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
		id := fmt.Sprint(i)
		be := id + strings.Repeat(" ", 20-len(id)) + "hello"

		if data != be {
			t.Error("data shoul be:", be)
		}
		t.Log(data, ".. ok")
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
