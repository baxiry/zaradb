package kvlite

import (
	"errors"
	"os"
	"testing"
)

func Test_Open(t *testing.T) {

	// test default dir db
	db := Open("")

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

	db.Close()
}

// main
func Test_Get_Put(t *testing.T) {

	db := Open("")
	defer db.Close()

	key := "hi"
	value := "hello world"

	db.Put(key, value)

	if db.Get(key) != value {
		t.Errorf("value shold be %s\n", value)
	}

}
