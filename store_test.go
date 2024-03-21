package store

import (
	"os"
	"path/filepath"
	"testing"
)

const maxItems = 11

var db *Database
var coll *Collection

func Test_NewDatabase(t *testing.T) {
	var path = "test"
	//ok
	db = NewDatabase(path)

	coll, _ = db.Collections["test"]

	_, err := os.Stat(filepath.Join(db.path, "test", "00000000000000000001"))
	if err != nil {
		if !os.IsExist(err) {
			t.Errorf("%s  should be exists", db.path+"/00000000000000000001")
		}
	}
}

func Test_NewCollection(t *testing.T) {
	if coll == nil {
		t.Error("coll should not be nil")
	}
}

func Test_reIndex(t *testing.T) {
	l := len(db.reIndex())
	if l != 2 {
		t.Error("lenght coll should be 2")
	}
}

func Test_insert(t *testing.T) {
	var i uint64

	data := "hello"

	l, _ := coll.log.LastIndex()
	for i = l; i <= l+maxItems; i++ {
		err := coll.insert(data)
		if err != nil {
			t.Log("normal error:", err)
		}
	}

	for i = l; i <= l+maxItems; i++ {
		data, err := coll.getData(i)
		if err != nil {
			t.Error(data, err)
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
		data, err := coll.getData(i)
		if err != nil {
			t.Errorf("no data witn %d idex\n", i)
		}

		if data != want {
			t.Errorf("shoul be: %s, not %s", want, data)
		}
	}
}

func Test_Close(t *testing.T) {
	indexs := coll.indexs
	if len(indexs) == 0 {
		t.Error("len of indexs should be great then 0")
	}

	if len(db.Collections)%maxItems-1 != 0 { // first item is 1 not 0
		t.Errorf("remain of 13 should be 0, not %d\n", len(db.Collections)%maxItems)
	}

	t.Log("indexs len: ", len(indexs))

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
