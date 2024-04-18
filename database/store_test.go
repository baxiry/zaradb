package database

import (
	"errors"
	"os"
	"testing"
)

func Test_Open(t *testing.T) {

	// test default dir db
	db := Open("")
	defer db.Close()

	if _, err := os.Stat("mok/"); errors.Is(err, os.ErrNotExist) {
		t.Errorf("err! %s should be exist", "mok/")
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

	db.Insert(coll, value)

	if db.Get(0, "") != value {
		t.Errorf("value shold be %s\n", value)
	}

}
