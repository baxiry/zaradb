package dblite

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/sjson"
)

// append new index in primary.index file
func NewIndex(index int64, dataSize int, indexFile *os.File) {
	strInt := fmt.Sprint(index) + " " + fmt.Sprint(dataSize)
	numSpaces := IndexLen - len(strInt)

	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	indexFile.WriteString(strInt)
}

// Insert
func Insert(path, data string) (err error) {

	PrimaryIndex++

	value, _ := sjson.Set(data, "_id", PrimaryIndex)
	println(value)

	path += fmt.Sprint(PrimaryIndex / 1000)
	fmt.Println("Path is ", path)
	_, err = Append(data, pages.Pages[path])

	// index
	NewIndex(PrimaryIndex, len(data), pages.Pages[indexFilePath])

	return err
}

// Select reads data form docs
func Select(path string) (data string) {

	return data
}

// Update update document data
func Update(serial, data string) (err error) {
	return
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

// append data to Pagefile & returns file size or error
func Append(data string, file *os.File) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// LastIndex return last index in table
func lastIndex(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size() / 20
}

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

	// out the buffer content
	return string(buffer[:n])
}

// index ingene

const IndexLen = 20

// get pageName Data Location  & data size from primary.indexes file
func GetIndex(id int, IndxFile *os.File) (pageName string, at, size int64) {

	pageName = strconv.Itoa(int(id) / 1000)
	bData := make([]byte, 20)
	_, err := IndxFile.ReadAt(bData, int64(id*20))
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

// deletes index from primary.index file
func DeleteIndex(id int, indxfile *os.File) { //
	at := int64(id * 20)
	indxfile.WriteAt([]byte("                    "), at)
}
