package dblite

import (
	"fmt"
	"io"
	"os"

	"github.com/Jeffail/gabs/v2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// data enginge

// At is where enginge insert data in page
var At int

// Select reads data form docs
func SelectById(query string) (result string) {
	id := gjson.Get(query, "where_id")

	if int(id.Int()) >= len(IndexsCache.indexs) {
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := IndexsCache.indexs[id.Int()][0]
	size := IndexsCache.indexs[id.Int()][1]

	// TODO fix page path
	result = Get(pages.Pages[RootPath+"0"], at, int(size))

	return result
}

// Update update document data
func Update(path, query string) (result string) {

	data := SelectById(query)
	fmt.Printf("DATA: %v\n", data)

	newData := gjson.Get(query, "data")
	fmt.Println("New data : ", newData)

	// `{"object":{"first":1,"second":2,"third":3}}`
	jsonParsed, err := gabs.ParseJSON([]byte(newData.String()))
	if err != nil {
		return fmt.Sprintf("ERROR: parse data json %s", err)
	}

	// extract fields that need to update
	for field, val := range jsonParsed.ChildrenMap() {

		result, _ = sjson.Set(data, field, val)

		data = result

		fmt.Printf("DATA: %v\n", data)
	}

	//path += fmt.Sprint(PrimaryIndex / 1000)

	err = InsertData(path, data)
	if err != nil {
		return "Fielure Insert"
	}

	//id := gjson.Get(query, "_id")
	//UpdateIndex(int(id.Int()), int64(at), int64(len(result)), pages.Pages[indexFilePath])

	fmt.Printf("updated data : %v\n", data)
	return "Success update"
}

// InsertData isert data directly (wethout extract data from query)
func InsertData(path, data string) (err error) {

	id := gjson.Get(data, "_id")

	path += fmt.Sprint(PrimaryIndex / 1000)

	at, err := Append(data, pages.Pages[path])
	if err != nil {

		return fmt.Errorf("ERROR! from Append %v\n", err)
	}

	// Update index
	size := int64(len(data))

	fmt.Println("We are in Insert Data")
	fmt.Println("id : ", id, "at : ", at, "size : ", size)

	UpdateIndex(int(id.Int()), int64(At), size, pages.Pages[indexFilePath])

	At += int(size)

	return err
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

	// TODO update index in indexsCache
	fmt.Println("IndexCace befor\n", IndexsCache.indexs)
	IndexsCache.indexs[id] = [2]int64{indexData, size}
	fmt.Println("IndexCache after: \n", IndexsCache.indexs)
}

// Insert
func Insert(path, query string) (res string) {

	data := gjson.Get(query, "data")

	value, err := sjson.Set(data.String(), "_id", PrimaryIndex)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}
	PrimaryIndex++

	path += fmt.Sprint(PrimaryIndex / 1000)

	size, err := Append(value, pages.Pages[path])
	if err != nil {
		fmt.Println("Error when append is : ", err)
		return "Fielure Insert"
	}

	// index
	NewIndex(At, len(value), pages.Pages[indexFilePath])
	At += size

	return fmt.Sprintf("Success Insert. _id : %d\n", PrimaryIndex-1)
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

// appends data to Pagefile & returns file size or error
func Append(data string, file *os.File) (size int, err error) {
	size, err = file.WriteString(data)
	if err != nil {
		println("Error WriteString ", err)
	}
	return size, err
}

// Delete removes document
func Delete(path string) (err error) {
	return
}
