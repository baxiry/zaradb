package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// convertAt convert location string to at and size int64
func convertAt(location string) (at, size int64) {

	// id here convert to at
	loc := strings.Trim(location, " ")
	sloc := strings.Split(loc, "-")

	id, _ := strconv.Atoi(sloc[0])

	siz, _ := strconv.Atoi(sloc[1])

	return int64(id), int64(siz)
}

func main() {

	at, size := convertAt("1755-666666")
	fmt.Println(at)
	fmt.Println(size)
	os.Exit(0)

	dbFile := "../example.db"

	file, err := Opendbs(dbFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data := "01234567890123456789"
	AppendData(file, data)

	getedData := GetVal(file, 0, 14)

	fmt.Println(getedData)
}

// TODO check if named returns improves performence ?!

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
