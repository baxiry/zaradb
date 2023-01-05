package main

import "time"

func main() {
	length := 1000000
	table := make(map[int]string, length)
	for i := 0; i < length; i++ {
		table[i] = "hello how ar you"
	}
	time.Sleep(time.Second * 7)
	table = nil
	time.Sleep(time.Second * 7)
}
