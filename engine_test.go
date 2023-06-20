package main

import (
	"fmt"
	"testing"
)

func Test_GetIndex(t *testing.T) {
	id := 111222
	page, at := GetIndex(id)
	if at != 222*IndexLen {
		t.Fatal("at must be 111 not", at)
	}
	if page != "111" {
		t.Fatal("page must be 222 not", at)
	}
	fmt.Println("GetIndex Done")
}

func Test_ConvIndex(t *testing.T) {
	location := "111 222   "
	at, size := ConvIndex(location)
	if at != 111 {
		t.Fatal("at must ber 111 not", at)
	}
	if size != 222 {
		t.Fatal("size must ber 222 not", size)
	}
	fmt.Println("ConvIndex Done")
}
