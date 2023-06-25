package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_UpdateIndex(t *testing.T) {
}

func Test_NewIndex(t *testing.T) {
	file, _ := os.OpenFile("primary.indexs", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	defer file.Close()

	for i := 1; i < 1002; i++ {
		NewIndex(i, file)
	}

	pageName, indx := GetIndex(14)
	fmt.Println(pageName, indx)

	pageName, indx = GetIndex(140)
	fmt.Println(pageName, indx)

	fmt.Println("NewIndex func Done")

}

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
