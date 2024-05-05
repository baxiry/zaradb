package engine

import (
	"errors"
	"os"
	"testing"
)

func Test_Open(t *testing.T) {

	// test default dir db
	db := Open("tmp")
	defer db.Close()

	if _, err := os.Stat("tmp/"); errors.Is(err, os.ErrNotExist) {
		t.Errorf("err! %s should be exist", "tmp/")
	}

	err := os.RemoveAll("mok/")
	check("", err)

	// test named dir db
	dir := "mydb/"

	db = Open(dir)
	db.Close()

	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		t.Errorf("err! %s should be exist", dir)
	}

	err = os.RemoveAll(dir)
	check("", err)

}

// main
func Test_Get_Put(t *testing.T) {

	db := Open("")
	defer db.Close()

	coll := "users"
	value := "hello world"

	//db.insert(coll, value)
	_, _ = coll, value

	if db.Get(0, "") != value {
		t.Errorf("value shold be %s\n", value)
	}

}
