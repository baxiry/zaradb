package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func ListDir(path string) {
	dbs, err := os.ReadDir(rootPath + path)
	if err != nil {
		fmt.Println(err)
	}

	dirs := 0
	for _, dir := range dbs {
		if dir.IsDir() && string(dir.Name()[0]) != "." {
			dirs++
			print(dir.Name(), " ")
		}
	}
	if dirs > 0 {
		println()
		return
	}
	println(path, "is impty")
}

// seedRand , use this func suparatly for nice performence
func seedRand() {
	rand.Seed(time.Now().Unix())
}
