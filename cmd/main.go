package main

import (
	"fmt"
)

var max = 10

func main() {
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
