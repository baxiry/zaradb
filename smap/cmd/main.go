package main

import (
	"fmt"
	"time"
	"zaradb/smap"
)

var sm = smap.NewSmap()
var m = map[string]string{}

func main() {
	for i := 0; i <= 100; i++ {
		id := fmt.Sprint(i)

		k := "hi_" + id
		v := "hello_" + id

		sm.Set(k, v)
		m[k] = v
	}

	fmt.Println("\nsmap")
	s := time.Now()
	for j := 0; j < 1000000; j++ {

		v := sm.Get("hi_20")
		if v != "hello_20" {
			panic("what")
		}

	}

	fmt.Println(time.Since(s))
	fmt.Println("\nmap")
	s = time.Now()

	for j := 0; j < 1000000; j++ {

		v, _ := m["hi_20"]
		if v != "hello_20" {
			panic("what")
		}
	}
	fmt.Println(time.Since(s))
}
