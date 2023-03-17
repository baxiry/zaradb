package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// PageLength is count of items for each page
const PageLength = 1000

// TODO

// writeIndex
// readIndex

func main() {

	NewPage(999)

	file, _ := Opendbs("../example.db")
	defer file.Close()

	AppendData(file, "01234567890123456789 ")

	getedData := GetVal(file, 0, 14)

	fmt.Println("geted data:", getedData)
}

// convertAt convert  string location to at and size int64
func convertIndex(location string) (at, size int64) {

	sloc := strings.Split(location, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}

// getLocation take id and return pageName and indexLocation
func getLocation(id string) (string, int) {
	indx, _ := strconv.Atoi(id[len(id)-3:])
	return id[:len(id)-3], indx
}

// NewPage create new file db page
func NewPage(id int) {
	// page is a file with som headrs to store data

	id = id / PageLength
	fmt.Println("file name is ", id)

	sid := strconv.Itoa(id)

	file, err := os.Create(sid)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// fill header file
	for i := 0; i < 1000; i++ {
		if _, err := file.WriteString("          "); err != nil {
			panic(err)
		}
	}
}

// Opendb opens | create new file
func Opendbs(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

// GetVal returns data as string.
// it take file pointr, at int64 & len of data that will read
func GetVal(file *os.File, at int64, buff int) string {
	// TODO check if reusing global buffer fast !
	buffer := make([]byte, buff)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return ""
	}
	// out the buffer content
	return string(buffer[:n])
}

// AppendData appends data to file
// return lenght of file or err
func AppendData(file *os.File, data string) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// getField get field from json string
func getField(field, json string) string {
	return gjson.Get(json, field).String()
}

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last + 1
}
