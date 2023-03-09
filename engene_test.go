package main

import (
	"testing"
)

func Test_convertAt(t *testing.T) {
	location := "111-222   "
	at, size := convertAt(location)
	if at != 111 {
		t.Fatal("at must ber 111 not", at)
	}
	if size != 222 {
		t.Fatal("size must ber 222 not", at)
	}

}
