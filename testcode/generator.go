package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func main() {

	for i := 0; i < 1000; i++ {
		n := strconv.Itoa(i)
		d, err := ioutil.ReadFile("../../test/" + n + "_file.txt")
		if err != nil {
			panic(err)
		}
		// fmt.Println(string(d))
		_ = string(d)

	}
	os.Exit(1)

	s := "This is a string\n"

	write()

	for i := 0; i < 1000000; i++ {
		snum := strconv.Itoa(i)
		ioutil.WriteFile(snum+"_file.txt", []byte(snum+" "+s), 0644)
	}
}

func write() {

	//Write first line
	err := ioutil.WriteFile("temp.txt", []byte("first line\n"), 0644)
	if err != nil {
		fmt.Println(err)
	}

	//Append second line
	f, err := os.OpenFile("temp.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	for i := 0; i < 1000000; i++ {
		snum := strconv.Itoa(i)
		//write directly into file
		f.Write([]byte(snum + " This is a string\n"))

	}
}
