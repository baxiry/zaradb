package dblite

import (
	"io"
	"os"
)

// data enginge

// appends data to Pagefile & returns file size or error
func Append(file *os.File, data string) (size int, err error) {

	size, err = file.WriteAt([]byte(data), collect.at)
	if err != nil {
		eLog.Println("Error WriteString ", err)
	}
	return size, err
}

// gets data from *file, takes at (location) & buffer size
func Get(file *os.File, at int64, size int) string {

	buffer := make([]byte, size)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		eLog.Println(err)
		return "Error! check if collection name is correct"
	}

	// out the buffer content
	return string(buffer[:n])
}
