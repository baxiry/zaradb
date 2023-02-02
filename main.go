package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var data = `
%d
helle world ;sjd;fja sjafa sajsdfjasf j asjdfa sfoja faosjf;oajsf;o asfja;sfj asjf asjf;asj fsoa\n
`

func main() {

	start := time.Now()
	file1, err := os.OpenFile("file1.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	file2, err := os.OpenFile("file2.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	defer file1.Close()
	defer file2.Close()

	go func() {

		for i := 0; i < 1000000; i++ {
			file1.WriteString(fmt.Sprintf(data, i))
		}
	}()
	for i := 0; i < 1000000; i++ {
		file2.WriteString(fmt.Sprintf(data, i))
	}
	fmt.Println("write duration", time.Since(start))
	start = time.Now()

	var lendata int
	for i := 0; i < 1000000; i++ {
		lendata += len(GetVal(file1, int64(i), i+10))
	}

	fmt.Println(lendata)

	var lendata2 int
	for i := 0; i < 1000000; i++ {
		lendata2 += len(GetVal(file2, int64(i), i+10))
	}
	fmt.Println(lendata2)

	fmt.Println("read duration", time.Since(start))
	fmt.Println("Done")

	dbFile := "../example.db"

	file, err := Opendbs(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data := "tested data ok"
	AppendData(file, data)

	getedData := GetVal(file, 0, 14)

	fmt.Println(getedData)
}

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return strings.Split(args[1], ".")
}
