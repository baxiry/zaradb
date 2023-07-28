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
	id := gjson.Get(query, "whereId")

	for i, v := range indexsCache.indexs {
		fmt.Println("i is : ", i, "v is : ", v)
	}

	println("id in SelectId is : ", id.Int())

	at := indexsCache.indexs[id.Int()][0]

	println("at in SelectId is : ", at)

	size := indexsCache.indexs[id.Int()][1]

	println("size in SelectId is : ", size)

	result = Get(pages.Pages[RootPath+id.String()], at, int(size))

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
		return ""
	}

	// out the buffer content
	return string(buffer[:n])
}

var At int

// Insert
func Insert(path, data string) (res string) {

	value, _ := sjson.Set(data, "_id", PrimaryIndex)
	println(value)

	path += fmt.Sprint(PrimaryIndex / 1000)
	fmt.Println("Path is ", path)

	at, err := Append(data, pages.Pages[path])
	if err != nil {
		fmt.Println("Error when append is : ", err)
		return "Fielure Insert"
	}

	fmt.Println("at in insert func is :", at)

	// index
	NewIndex(At, len(data), pages.Pages[indexFilePath])

	At += at
	PrimaryIndex++

	return fmt.Sprintf("Success Insert. _id is %d", PrimaryIndex-1)
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
