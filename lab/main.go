package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	path := "example.data"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := "Hi_World! "
	file2, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file2.Close()

	start := time.Now()

	src := ""

	for i := 0; i < 1000000; i++ {
		// src = getVal(file, int64(i+34))
		AppendData(file2, data)
	}

	fmt.Println(src)
	fmt.Println("Done ", time.Since(start))

}

func AppendData(file *os.File, data string) {

	_, err := file.WriteString(data)
	if err != nil {
		fmt.Println(err)
	}
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

func LastSerial(path string) (str string, err error) {

	return "", nil
}

func findSerial(gool uint) (uint, error) {

	staps := 0
	base := uint(0)
	big := uint(1)
	smal := uint(0)

	for staps = 0; staps <= 300; staps++ {

		if base < gool {
			smal = base
			base = big
			big = base * 2
			if big >= (2305843009213693000) {
				// TODO find bigger then this number
				return 0, fmt.Errorf("big overflow at stap: %d\n", staps)
			}
		}

		if base > gool {
			big = base
			base = (smal + big) / 2
		}

		if base == gool {
			if base+1 > gool {
				println("gool", base)
				println("staps", staps)
				break
			}
		}
	}

	return base, nil
}

func genSerial() string {
	src := []string{"a", "b", "b", "c", "d", "e", "f", "j", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	res := ""

	for j := 0; j < 1000000; j++ {
		res = ""
		for i := 0; i < 5; i++ {
			res += string(src[i])
		}
	}

	fmt.Println("res : ", res)
	return res
}

func goLoop() {

	var mems = []string{"dog", "cat", "mouse", "koko", "bebe", "jojo", "haha", "jiji", "foo", "bar", "bax", "hik", "jik", "ors", "nos", "ren"}

	var bots = []string{"bot_1", "bot_2", "bot_3", "bot_4", "bot_5", "bot_6", "bot-7"}

	i := 0
	lbots := len(bots)
	for _, mem := range mems {
		go func(mem string, i int) {
			fmt.Printf("%s kik %s\n", bots[i], mem)
		}(mem, i)
		i++
		if lbots == i {
			i = 0
		}
	}
	time.Sleep(time.Millisecond * 10)
}
