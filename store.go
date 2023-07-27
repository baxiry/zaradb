package dblite

import (
	"fmt"
	"io"
	"os"

	"github.com/tidwall/sjson"
)

// data enginge

// Select reads data form docs
func Select(id int) (data string) {
	return data
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

// Update update document data
func Update(serial, data string) (err error) {
	return
}

// Delete removes document
func Delete(path string) (err error) {
	return
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
