package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//last serial is 10010000

var wg sync.WaitGroup

func main() {

	start := time.Now()

	args := os.Args
	if len(args) < 2 {
		fmt.Println("inter arg")
		return
	}

	loops, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// read data

	lendata := 0

	ch := make(chan string, 2)
	var str string

	for i := 0; i < loops; i++ {

		v := rand.Intn(201000-1) + 1 // range is min to max

		wg.Add(1)
		go func() {
			d, _ := ioutil.ReadFile("/Users/fedora/repo/test/" + strconv.Itoa(v) + "_file.txt")
			s := string(d[0])
			ch <- s
		}()
	}

	for i := 0; i < loops; i++ {
		str += <-ch
	}

	lendata += len(str)

	println("size of data : ", lendata)
	println("duration ms  : ", time.Since(start).Milliseconds())

	os.Exit(0)

	data := args[2]
	for i := 0; i < loops; i++ {

		n := strconv.Itoa(i)
		path := "/Users/fedora/repo/test/" + n + "_file.txt"
		err = os.WriteFile(path, []byte(n+data+"\n"), 0644)
		if err != nil {
			println("err is : ", err)
		}

	}

	wg.Wait()
	println("duration: ", time.Since(start).Milliseconds())
	os.Exit(0)

	////////////////////////////////

}

// create new files with new data
func creatFiles(from, to int) (err error) {
	for i := from; i < to; i++ {
		n := strconv.Itoa(i)
		path := "/Users/fedora/repo/test/" + n + "_file.txt"
		err = os.WriteFile(path, []byte(n+" this is a stirng"), 0644)
		if err != nil {

			return err
		}
	}
	return nil
}

// generate data to file
func generateLinse(n int) {
	s := "This is a string\n"

	for i := 0; i <= n; i++ {
		snum := strconv.Itoa(i)
		ioutil.WriteFile(snum+"_file.txt", []byte(snum+" "+s), 0644)
	}

}

func generateLinesToFile(n int, file string) {
	//Append second line
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
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

// read multi files
func readFiles(from, to int) {

	for i := from; i < to; i++ {
		n := strconv.Itoa(i)
		d, err := ioutil.ReadFile("/Users/fedora/repo/test/" + n + "_file.txt")
		if err != nil {
			panic(err)
		}
		fmt.Print(string(d))
	}
}

func readSpicificLine(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}
