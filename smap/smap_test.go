package smap

import (
	"fmt"
	"testing"
	"time"
)

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}

func BenchmarkMap(t *testing.B) {

	keys := list{}

	mkeys := map[string]string{}

	for i := 0; i < 20; i++ {
		v := "hello_" + fmt.Sprint(i)
		//fmt.Println("v is : ", v)
		keys.set(v, v)
		mkeys[v] = v
	}

	s := time.Now()
	tot := 0

	for i := 0; i < max; i++ {
		if keys.get("hello_0") == "hello_0" {
			tot++
		}
	}

	fmt.Println(tot)
	fmt.Println(time.Since(s))
	fmt.Println("my len: ", keys.len())

	s = time.Now()
	tot = 0

	for i := 0; i < max; i++ {
		if mkeys["hello_0"] == "hello_0" {
			tot++
		}
	}

	fmt.Println(tot)
	fmt.Println(time.Since(s))

	fmt.Println("map len: ", len(mkeys))
}
