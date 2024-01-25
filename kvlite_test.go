package kvlite

import (
	"testing"
)

func Test_Open(t *testing.T) {
}

// main
func Test_insert_get(t *testing.T) {

	db, err := NewWriter("test")
	if err != nil {
		t.Errorf("inicialaze error")
	}
	defer db.Close()

	value := "hello world"
	db.insert(value)
	v, err := db.get(1)
	if err != nil {
		t.Errorf("err when get data")
	}

	if v != value {
		t.Errorf("value shold be %s, not %s\n", value, v)
	}

}
