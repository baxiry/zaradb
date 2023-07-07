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
func NewIndex(ind, dsize int, file *os.File) { // dsize is data size
	strInt := fmt.Sprint(ind) + " " + fmt.Sprint(dsize)
	numSpaces := IndexLen - len(strInt)

	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	file.WriteString(strInt)
}

// get pageName  indexLocation  & data size from primary.indexes file
func GetIndex(id int, indxFile *os.File) (pageName string, at, size int64) {

	pageName = strconv.Itoa(int(id) / 1000)
	at = int64(id % 1000)

	bData := make([]byte, 20)
	_, err := indxFile.ReadAt(bData, at)
	if err != nil {
		println(err)
		return
	}

	sData := strings.Split(string(bData), " ")[1]

	isize, _ := strconv.Atoi(fmt.Sprint(sData))

	return pageName, at * IndexLen, int64(isize)
}

// update index val in primary.index file
func UpdateIndex(ind int64, file *os.File) {

	file.WriteString(fmt.Sprint(ind))
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
