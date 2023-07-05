package main

import (
	"os"
	"testing"
)

func Test_UpdateIndex(t *testing.T) {
}

func Test_NewIndex(t *testing.T) {
	file, _ := os.OpenFile("primary.indexs", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	defer func() {
		file.Close()
		os.Remove("primary.indexs")
	}()

	for i := 0; i < 1111; i++ {
		NewIndex(i, file)
	}

	// this func for get exact endix location in file "At"
	res := func(i int64) int64 { return (i % 1000) * 20 }

	// input 1
	pageName, indx := GetIndex(1)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != res(1) { // (1 % 1000) * 20 = 20
		t.Error("index must be 20")
	}

	//"input 140 return 2800
	pageName, indx = GetIndex(140)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != res(140) { // 2800
		t.Error("index must be 2800")
	}

	//"input 1111: 2220
	pageName, indx = GetIndex(1111)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}

	if indx != res(1111) {
		t.Error("index must be 2220")
	}
	println("NewIndex \t Done")

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
	println("GetIndex \t Done")
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

	println("ConvIndex \t Done")
}
