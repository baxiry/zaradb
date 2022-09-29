package main

import (
	"math/rand"
	"sync"
	"time"
)

var (
	wg       sync.WaitGroup
	Latters  = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"} // "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	ListLen  = len(Latters) - 2
	rootPath = "" //"/Users/fedora/.mydb/test/"
)

func init() {
	seedRand()
}

func seedRand() {
	rand.Seed(time.Now().Unix())
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5)
	//		rand.Seed(time.Now().Unix())
	//	}
	//}()
}
