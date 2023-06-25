package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const IndexLen = 20

const IndexsFile = "primary.index"

// append new index in primary.index file
func NewIndex(ind int, file *os.File) {
	strInt := fmt.Sprint(ind)
	nSpaces := IndexLen - len(strInt)
	for i := 0; i < nSpaces; i++ {
		strInt += " "
	}

	file.WriteString(strInt)
}

// update index val in primary.index file
func UpdateIndex(ind int64, file *os.File) {

	file.WriteString(fmt.Sprint(ind))
}

// get pageName & indexLocation where data is stored
func GetIndex(id int) (pageName string, at int64) {
	pageName = strconv.Itoa(int(id) / 1000)
	at = int64(id % 1000)
	return pageName, at * IndexLen
}

// get index data & pageName from primary.index
func GetWhere(id int) (pageName string, at int64) {
	// should read primary.index here
	return pageName, at * IndexLen
}

// converts index location to at and size as int64
func ConvIndex(IndexLocation string) (at, size int64) {
	sloc := strings.Split(IndexLocation, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}

// gets data from *file, takes at (location) & buffer size
func Get(file *os.File, at int64, size int) string {

	buffer := make([]byte, size)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("readAt ", err)
		return ""
	}
	// out the buffer content
	return string(buffer[:n])
}

// append data to Pagefile & returns file size or error
func Append(file *os.File, data string) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// gets field from json
func getField(field, json string) string {
	return gjson.Get(json, field).String()
}

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last + 1
}
