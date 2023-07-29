package dblite

import (
	"fmt"
	"io"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// data enginge

// Select reads data form docs
func SelectId(query string) (result string) {
	id := gjson.Get(query, "where_id")

	if int(id.Int()) >= len(indexsCache.indexs) {
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := indexsCache.indexs[id.Int()][0]
	size := indexsCache.indexs[id.Int()][1]

	// TODO fix page path
	result = Get(pages.Pages[RootPath+"0"], at, int(size))

	return result
}

// Select reads data form docs
func Select(filter string) (result string) {
	id := gjson.Get(filter, "_id")
	fmt.Println("id is ", id.String())

	return result
}

// gets data from *file, takes at (location) & buffer size
func Get(file *os.File, at int64, size int) string {

	buffer := make([]byte, size)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("readAt ", err)
		return "ERROR form ReadAt func"
	}

	// out the buffer content
	return string(buffer[:n])
}

var At int

// Insert
func Insert(path, query string) (res string) {

	data := gjson.Get(query, "data")

	value, err := sjson.Set(data.String(), "_id", PrimaryIndex)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	path += fmt.Sprint(PrimaryIndex / 1000)
	fmt.Println("Path is ", path)

	at, err := Append(value, pages.Pages[path])
	if err != nil {
		fmt.Println("Error when append is : ", err)
		return "Fielure Insert"
	}

	// index
	NewIndex(At, len(value), pages.Pages[indexFilePath])

	At += at
	PrimaryIndex++

	return fmt.Sprintf("Success Insert. _id : %d\n", PrimaryIndex-1)
}

// append data to Pagefile & returns file size or error
func Append(data string, file *os.File) (int, error) {
	fileSize, err := file.WriteString(data)
	if err != nil {
		println("Error WriteString ", err)
	}
	return fileSize, err
}

// Update update document data
func Update(serial, data string) (err error) {
	return
}

// Delete removes document
func Delete(path string) (err error) {
	return
}
