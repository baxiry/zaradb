package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(GenSerial(5))
	}
}

// GenSerial generate serial for Doc
func GenSerial(length int) (serial string) {
	var i int
	for i = 0; i < length; i++ {
		serial += Latters[rand.Intn(ListLen)+1]
	}
	return serial
}

// Delete remove document
func Delete(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return
}

// create new files with new data
func Insert(path, data string) (err error) {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("inter numver arg")
		return
	}

	err = os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return
	}
	return
}

// Select reads data form docs
func Select(path string) (data string, err error) {

	bdata, err := os.ReadFile(rootBase + path)
	if err != nil {
		return "", err
	}
	return string(bdata), nil

}
