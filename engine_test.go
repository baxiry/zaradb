package main

import (
	"fmt"
	"testing"
)

func Test_getAt(t *testing.T) {
	id := 111222
	page, at := GetAt(id)
	if at != 222*LenIndex {
		t.Fatal("at must be 111 not", at)
	}
	if page != "111" {
		t.Fatal("page must be 222 not", at)
	}
	fmt.Println("getAt Done")
}

func Test_convIndex(t *testing.T) {
	location := "111 222   "
	at, size := ConvIndex(location)
	if at != 111 {
		t.Fatal("at must ber 111 not", at)
	}
	if size != 222 {
		t.Fatal("size must ber 222 not", size)
	}
	fmt.Println("convertIndex Done")
}
