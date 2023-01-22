package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	appendData("temp.txt")
	cleanFile("temp.txt")
}

func cleanFile(fname string) {
	err := ioutil.WriteFile(fname, []byte(""), 0644)
	if err != nil {
		log.Fatal(err)
	}

}
func appendData(fname string) {

	//Append second line
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	for i := 0; i < 10; i++ {
		numberLine := strconv.Itoa(i)
		if _, err := file.WriteString(numberLine + "new data\n"); err != nil {
			log.Fatal(err)
		}
	}
	//Print the contents of the file
	data, err := ioutil.ReadFile("temp.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

}
