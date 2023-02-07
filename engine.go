package main

import (
	"fmt"
	"io"
	"os"

	"github.com/tidwall/gjson"
)

// Opendb opens | create  new file db
func Opendbs(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

// GetVal return data string.
// take file pointr, at int64 & len of data will read
func GetVal(file *os.File, at int64, buff int) string {
	// TODO my be nice idea to make buffer global  to reuse it
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

// AppendData appends data string to file
// return len or size of file and err
func AppendData(file *os.File, data string) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// getField get field from json string
func getField(field, json string) string {
	return gjson.Get(json, field).String()
}
