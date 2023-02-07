package main

import "fmt"

// simplest query language
func queryLang() {
	query := arguments()
	fmt.Println("query is : ", query)
}

// Update update document data
func Insert(path, data string) (err error) {
	// TODO add ''where'' statment insteade by serial
	return
}

// rootPath = "/Users/fedora/.mydb/test/"

// Select reads data form docs
func Select(path string) (data string) {
	return data
}

// Update update document data
func Update(serial, data string) (err error) {
	// TODO add ''where'' statment ensteade serial
	return
}

// Delete remove document
func Delete(path string) (err error) {
	return
}

// get json data from stdin argument
func getQueryJsonArgs(str string) (json string) {
	var start, end int32

	var i int32
	for i = 0; i < int32(len(str)); i++ {
		if str[i] == '{' {
			start = i
			break
		}
	}
	for i = int32(len(str)) - 1; i >= 0; i-- {
		if str[i] == '}' {
			end = i
			break
		}
	}
	return str[start : end+1]
}
