package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// data enginge

// gets data from *file, takes at (location) & buffer size
func Get(file *os.File, at int64, size int) string {

	buffer := make([]byte, size)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("readAt ", err)
		return ""
	}
	println("data is : ", string(buffer[:n]))

	// out the buffer content
	return string(buffer[:n])
}

// append data to Pagefile & returns file size or error
func Append(data string, file *os.File) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last
}

// index ingene

const IndexLen = 20

const IndexsFile = "primary.index"

// get pageName  Data Location  & data size from primary.indexes file
func GetIndex(id int, indxFile *os.File) (pageName string, at, size int64) {

	pageName = strconv.Itoa(int(id) / 1000)
	bData := make([]byte, 20)
	_, err := indxFile.ReadAt(bData, int64(id*20))
	if err != nil {
		panic(err)
	}

	slc := strings.Split(string(bData), " ")
	iat, _ := strconv.Atoi(slc[0])

	isize, _ := strconv.Atoi(fmt.Sprint(slc[1]))

	return pageName, int64(iat), int64(isize)
}

// update index val in primary.index file
func UpdateIndex(id int, indexData, size int64, indexFile *os.File) {

	at := int64(id) * 20

	strIndex := fmt.Sprint(indexData) + " " + fmt.Sprint(size)
	for i := len(strIndex); i < 20; i++ {
		strIndex += " "
	}

	_, err := indexFile.WriteAt([]byte(strIndex), at)
	if err != nil {
		panic(err)
	}

}

// append new index in primary.index file
func NewIndex(ind, dsize int, file *os.File) { // dsize is data size
	strInt := fmt.Sprint(ind) + " " + fmt.Sprint(dsize)
	numSpaces := IndexLen - len(strInt)

	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	file.WriteString(strInt)
}

// deletes index from primary.index file
func DeleteIndex(id int, indxfile *os.File) { //
	at := int64(id * 20)
	indxfile.WriteAt([]byte("                    "), at)
}
