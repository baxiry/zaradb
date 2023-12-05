package main

import (
	"fmt"
	"kvlite"
	"time"
)

var max = 100_000

// main
func main() {

	// first db
	db := kvlite.Open("db1/")
	defer db.Close()

	s := time.Now()
	for i := 0; i < max; i++ {
		db.Set(i, "users", "hello world:"+fmt.Sprint(i))
	}
	fmt.Println(time.Since(s))
	s = time.Now()

	// set data

	l := 0
	for i := 0; i < max; i++ {
		l += len(db.Get(i))
	}
	fmt.Println(time.Since(s))

	data := db.Get(333)
	fmt.Println("len & data:", l, data)

}
