package main

import (
	"os"
	"testing"
)

func Test_UpdateIndex(t *testing.T) {
}

func Test_All_Index_Func(t *testing.T) {
	file, _ := os.OpenFile("primary.indexs", os.O_RDWR|os.O_CREATE, 0644)
	defer func() {
		file.Close()
		//os.Remove("primary.indexs")
	}()

	// testing NewIndex func
	for i := 0; i <= 1111; i++ {
		NewIndex(i, i, file)
	}

	// this func for get exact endix location in file "At"
	res := func(i int64) int64 { return (i % 1000) * 20 }

	// testing GetIndex func
	// input 1
	pageName, indx, size := GetIndex(1, file)
	if pageName != "0" {
		t.Error("pageName must be 0")
	}
	if indx != res(1) { // (1 % 1000) * 20 = 20
		t.Error("index must be 20")
	}
	if size != 1 {
		t.Errorf("size is %d, shoul be %d ", size, 1)
	}

	//"input 140 return 2800
	pageName, indx, size = GetIndex(140, file)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != res(140) { // 2800
		t.Error("index must be 2800")
	}
	if size != 140 { // 2800
		t.Errorf("size is %d, must be %d", size, 140)
	}

	//"input 1111: 2220
	pageName, indx, size = GetIndex(1111, file)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}
	if indx != res(1111) {
		t.Errorf("index s %d, must be %d\n", indx, 2220)
	}
	if size != 1111 {
		t.Error("size must be ", 1111)
	}

	// testing UpdateIndex func
	for i := 10; i <= 1111; i++ {
		UpdateIndex(i, int64(i+5), int64(i+10), file)
	}

}
