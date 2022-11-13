package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	path := "example.data"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {

		println(err)
	}
	defer file.Close()

	file2, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	data := "Hi__World! "

	start := time.Now()

	src := ""

	for i := 0; i < 10; i++ {
		AppendData(file, data+strconv.Itoa(i))
		//src = getVal(file2, int64(i))
	}

	fmt.Println(src)
	fmt.Println("Done ", time.Since(start))
	fmt.Println("Zise", GetFileSize(path))
}

// AppendData to file opend with os.O_APPEND|os.O_CREATE|os.O_WRONLY params
func AppendData(file *os.File, data string) (err error) {
	a, err := file.WriteString(data)

	println("a is : ", a)
	fmt.Println("size : ")
	return
}

func getVal(file *os.File, at int64) string {
	buffer := make([]byte, 10)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		panic(err)
	}
	// out the buffer content
	return string(buffer[:n])
}

func GetFileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return file.Size()
}

func IsExest(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
