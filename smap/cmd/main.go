package main

import (
	"fmt"
	"time"
	"zaradb/smap"
)

var max = 1000000

var keys = smap.NewSmap()
var m = map[string]string{}

func main() {

	m := make(map[string]string)

	for i := 0; i <= 20; i++ {
		k := "hi_" + fmt.Sprint(i)
		v := "hello_" + fmt.Sprint(i)
		keys.Set(k, v)
		m[k] = v
	}

	s := time.Now()
	for i := 0; i < max; i++ {
		if keys.Get("hi_10") != "hello_10" {
			panic("hat")
		}
	}

	fmt.Println("smap", time.Since(s))

	st := time.Now()
	for i := 0; i < max; i++ {
		if m["hi_20"] != "hello_20" {
			panic("hat")
		}
	}

	fmt.Println("map", time.Since(st))

}
