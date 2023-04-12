package main

import (
	"db"
	"fmt"
	"os"
	"time"
)

func init() {
	fmt.Println("root path is ", db.RootPath)
}

func main() {

	page, _ := os.Create("example.db")
	defer page.Close()

	start := time.Now()

	byt := make([]byte, 1000)

	total := 0

	for i := 0; i < 1000000; i++ {
		page.ReadAt(byt, int64(i))

		total += len(byt)
	}

	fmt.Println("data size ", total)
	fmt.Println("Douration: ", time.Since(start))
	time.Sleep(time.Second * 15)

}
