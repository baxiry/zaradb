package main

import (
	"fmt"
	"zaradb/store"
)

func main() {
	eng := store.NewDatabase("mydb")
	for coll, ok := range eng.Collections {

		println(coll, ok)
	}

	// insert
	for i := 1; i < 12; i++ {
		//	store.Collections["test"].insert("hello_" + fmt.Sprint(i))
	}

	eng.Close()

	fmt.Println("path : ", store.Slash())

}
