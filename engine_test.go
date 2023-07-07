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
	defer func() {
		file.Close()
		//os.Remove("primary.indexs")
	}()

	for i := 0; i < 1111; i++ {
		NewIndex(i, 1024, file)
	}

	// this func for get exact endix location in file "At"
	res := func(i int64) int64 { return (i % 1000) * 20 }

	// input 1
	pageName, indx, size := GetIndex(1, file)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != res(1) { // (1 % 1000) * 20 = 20
		t.Error("index must be 20")
	}

	if size != 1024 {
		t.Error("size shoul be 1024")
	}
	fmt.Println("Size is ", size)

	//"input 140 return 2800
	pageName, indx, _ = GetIndex(140, file)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != res(140) { // 2800
		t.Error("index must be 2800")
	}

	//"input 1111: 2220
	pageName, indx, _ = GetIndex(1111, file)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}

	if indx != res(1111) {
		t.Error("index must be 2220")
	}
	println("NewIndex \t Done")

}
