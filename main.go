package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// PageLength is count of items for each page
const HeaderLength = 1000

type Pages struct {
	Pages map[string]*os.File
}

// writeIndex
// readIndex

func main() {

	var db Pages
	db.Open("../dbs")
	defer db.Close()

	//AppendData(file, "01234567890123456789 ")

	getedData := GetValue("../dbs/0", 0, 14)

	fmt.Println("geted data:", getedData)
}

// Opendb opnens all pages in root db
func (db *Pages) Open(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("readDir: ", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		page, err := os.OpenFile(file.Name(), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("os open file: ", err)
		}
		db.Pages[file.Name()] = page

		fmt.Println("file", file.Name())
	}
}

// Close All pages
func (db Pages) Close() {
	for _, p := range db.Pages {
		p.Close()
	}
}

// NewPage create new file db page
func NewPage(id int) {

	// page is a file with som headrs to store data
	id = id / HeaderLength
	fmt.Println("file name is ", id)

	sid := strconv.Itoa(id)

	file, err := os.Create(sid)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// fill header file
	for i := 0; i < 1000; i++ {
		if _, err = file.WriteString("          "); err != nil {
			fmt.Println("ok", err)
		}
	}
}

// getLocation take str 'id' and returns str 'pageName' & index Location
func getLocation(id string) (string, int) {
	indxL, _ := strconv.Atoi(id[len(id)-3:]) // ret pageName
	at := id[:len(id)-3]                     // ret IndexLocation
	return at, indxL
}

// convertIndex convert string index location to at and size int64
func convIndex(IndexLocation string) (at, size int64) {

	sloc := strings.Split(IndexLocation, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}

// GetVal returns data as string.
// it take file pointr, at int64 & len of data that will read
func GetValue(nfile string, at int64, buff int) string {
	// TODO check if reusing global buffer fast !
	buffer := make([]byte, buff)

	// read at
	var p Pages
	n, err := p.Pages[nfile].ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("readAt ", err)
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
