package main

import (
	"fmt"
	"kvlite"
	"time"
)

var max = 10

func cmain() {
	var i uint64

	w, err := NewWriter("mylog")
	if err != nil {
		panic(err)
	}

	for i = 1; i < 10; i++ {
		_, err := w.insert("hello world" + fmt.Sprint(i))
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for i = 1; i <= w.index; i++ {
		d, err := w.get(i)
		if err != nil {
			println(err.Error())
		}
		println(string(d))
	}
}

// 4_294_967_295
// main
func main() {

	db := kvlite.Open("db1/")
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
