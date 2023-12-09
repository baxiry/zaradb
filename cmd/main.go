package main

import (
	"fmt"
	"kvlite"
)

var max = 10

// main
func main() {

	// first db
	db := kvlite.Open("db1/")
	defer db.Close()

	lid := db.Lid
	println("insert")
	for i := 0; i < max; i++ {
		db.Insert("users", "hello world:"+fmt.Sprint(i+lid))
	}

	for i := 0; i < max; i++ {
		fmt.Println("get: ", db.Get(i))
	}

	println("delete")
	for i := 0; i < 7; i++ {
		//		db.Delete(i, "users")
	}

	for i := 0; i < max+lid; i++ {
		fmt.Println(i, "get: ", db.Get(i))
	}

}
