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

// "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
var Latters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

//last serial is 10010000

var wg sync.WaitGroup

// get last serial
func GenerateID(latters []string) (serial string) {
	var i int
	for i = 0; i < 10; i++ {
		serial += latters[rand.Intn(35)+1]
	}
	return serial
}

func main() {
	rand.Seed(time.Now().UnixNano())
	r := GenerateID(Latters)
	fmt.Println(r)

	for i := 0; i < 1000000; i++ {
		GenerateID(Latters)
	}
	//readsfiles()
	//fmt.Println(creatFiles())
}

// create new files with new data
func creatFiles() (err error) {
	start := time.Now()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("inter numver arg")
		return
	}

	to, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < to; i++ {
		n := strconv.Itoa(i)
		path := "/Users/fedora/repo/test/" + n + "_file.txt"
		err = os.WriteFile(path, []byte(n+" this is a stirng"), 0644)
		if err != nil {
			return
		}
	}
	println("duration ms  : ", time.Since(start).Milliseconds())
	return
}

// read data form Mutli files
func readsfiles() {
	start := time.Now()
	args := os.Args
	if len(args) < 2 {
		fmt.Println("inter numver arg")
		return
	}

	loops, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	lendata := 0

	for j := 0; j < 100; j++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < loops; i++ {
				//v := rand.Intn(201000-1) + 1 // range is min to max
				d, _ := os.ReadFile("/Users/fedora/repo/test/" + strconv.Itoa(i) + "_file.txt")
				lendata += len(string(d[0]))
			}
		}()

		wg.Wait()
	}

	println("size of data : ", lendata)
	println("duration ms  : ", time.Since(start).Milliseconds())
}

// find a line
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

// read
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

// generate data to file
func generateLinse(n int) {
	s := "This is a string\n"

	for i := 0; i <= n; i++ {
		snum := strconv.Itoa(i)
		ioutil.WriteFile(snum+"_file.txt", []byte(snum+" "+s), 0644)
	}

}
