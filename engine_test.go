package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_UpdateIndex(t *testing.T) {
}

var DataFile, _ = os.OpenFile("data.page", os.O_RDWR|os.O_CREATE, 0644)

var IndexFile, _ = os.OpenFile("primary.indexs", os.O_RDWR|os.O_CREATE, 0644)

// testing all data storeing functions
func Test_Append(t *testing.T) {
	var at int64
	for i := 0; i < 13; i++ {

		data := "hello world ok "
		data += fmt.Sprint(i)

		lenByte := int64(len(data))
		lb, err := Append(data, DataFile)

		if err != nil {
			fmt.Println("error is : ", err)
		}

		myData := Get(DataFile, at, lb)
		fmt.Printf("Data is %s: \nlen byte is %d\nlb is %d\n ", myData, lenByte, lb)
		at += int64(lb)
	}
}

// testing all index functions
func Test_All_Index_Funcs(t *testing.T) {

	// testing NewIndex func
	for i := 0; i <= 1111; i++ {
		NewIndex(i, i, IndexFile)
	}

	// testing GetIndex func
	//"input 140 return 2800
	pageName, indx, size := GetIndex(140, IndexFile)
	if pageName != "0" {
		t.Error("pageName must be 1")
	}
	if indx != 140 { // 2800
		t.Error("index must be 2800")
	}
	if size != 140 { // 2800
		t.Errorf("size is %d, must be %d", size, 140)
	}

	//"input 1111: 2220
	pageName, indx, size = GetIndex(1111, IndexFile)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}
	if indx != 1111 {
		t.Errorf("index s %d, must be %d\n", indx, 1111)
	}
	if size != 1111 {
		t.Error("size must be ", 1111)
	}

	// testing UpdateIndex func
	for i := 10; i <= 1111; i++ {
		UpdateIndex(i, int64(i+5), int64(i+10), IndexFile)
	}

	//"input 1111: 2220
	pageName, indx, size = GetIndex(1111, IndexFile)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}

	if indx != 1111+5 {
		t.Errorf("index s %d, must be %d\n", indx, 1111+5)
	}
	if size != 1121 {
		t.Error("size must be ", 1121)
	}

	// testing DeleteIndex func
	DeleteIndex(1091, IndexFile)

	pageName, indx, size = GetIndex(1091, IndexFile)
	if pageName != "1" {
		t.Error("pageName must be 1")
	}

	if indx != 0 {
		t.Errorf("index s %d, must be %d\n", indx, 0)
	}
	if size != 0 {
		t.Error("size must be ", size)
	}

}

func Test_finish(t *testing.T) {
	DataFile.Close()
	os.Remove("primary.indexs")

	IndexFile.Close()
	os.Remove("primary.indexs")
	println("Done")

}
