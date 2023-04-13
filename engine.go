package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

const LenIndex = 20
const IndexsFile = RootPath + "primary.index"

func StoreIndex(path string) {}

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last + 1
}

// getLocation return fileName & indexLocation where data stored
func GetLocation(id int) (fileName string, indexLocation int64) {
	fileName = strconv.Itoa(id / 1000)
	indexLocation = int64(id % 1000)
	return fileName, indexLocation * LenIndex
}

// convertIndex convert string index location to at and size int64
func ConvIndex(IndexLocation string) (at, size int64) {
	sloc := strings.Split(IndexLocation, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}

// GetVal returns data as string.
// it take file pointr, at int64 & len of data that will read
func GetValue(file *os.File, at int64, buff int) string {
	// TODO check if reusing global buffer fast !
	buffer := make([]byte, buff)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("readAt ", err)
		return ""
	}
	// out the buffer content
	return string(buffer[:n])
}

// AppendData appends data to file
// return file size & err
func AppendData(file *os.File, data string) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// getField get field from json string
func getField(field, json string) string {
	return gjson.Get(json, field).String()
}
