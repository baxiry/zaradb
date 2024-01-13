package main

import (
	"fmt"
	"time"
	database "zaradb/db"
)

var max = 10

// 4_294_967_295
// main
func main() {

	db := database.Open("db1/")
	defer db.Close()

	s := time.Now()
	for i := 0; i < max; i++ {
		db.Insert("users", " hello world:"+fmt.Sprint(i+1))
	}
	fmt.Print("insert result : ")
	fmt.Println(time.Since(s))

	s = time.Now()

	for i := 0; i < max; i++ {
		_ = db.Get(i, "users")
	}
	fmt.Print("get result : ")
	fmt.Println(time.Since(s))
	//time.Sleep(time.Second * 20)

	db.Delete(3, "users")
	fmt.Println("delete result : ")
	for i := 122; i < 124; i++ {
		fmt.Println(123, db.Get(123, "users"))
	}

	db.Update(3, "users", "new data")
	fmt.Println("update result : ")
	fmt.Println(3, db.Get(3, "users"))
}
