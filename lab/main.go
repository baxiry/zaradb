package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func genData(n int) (data string) {
	num := strconv.Itoa(n)
	data = num
	for i := 0; i < 10-len(num); i++ {
		data += "_"

	}
	return data
}

func main() {
	path := "example.data"
	file, err := Opendb(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	src := ""

	for i := 0; i < 1000; i++ {
		AppendData(file, genData(i))
	}

	fmt.Println("size", FileSize(path))

	for i := 0; i < 1000; i++ {
		src = getVal(file, int64(10*i), 10)
		fmt.Println(src)
	}

}

// AppendData to file
func AppendData(file *os.File, data string) (err error) {
	lnb, err := file.WriteString(data)
	println("len bytes is : ", lnb)
	return
}

func getVal(file *os.File, at int64, ln int) string {
	buffer := make([]byte, ln)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("file size is : ", FileSize(file.Name()))
		fmt.Println("at is ", at)
		panic(err)
	}
	// out the buffer content
	return string(buffer[:n])
}

func Opendb(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

func FileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return file.Size()
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
