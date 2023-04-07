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

// getLocation return fileName & indexLocation where data stored
func getLocation(id int) (fileName string, indexLocation int64) {
	fileName = strconv.Itoa(id / 1000)
	indexLocation = int64(id % 1000)
	return fileName, indexLocation * LenIndex
}

// convertIndex convert string index location to at and size int64
func convIndex(IndexLocation string) (at, size int64) {
	sloc := strings.Split(IndexLocation, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}

// generateIndex
// writeIndex
// readIndex

func main() {

	fname, at := getLocation(2000)
	fmt.Printf("page is : %s at : %d\n", fname, at)

	dbPath := "/Users/fedora/repo/dbs/"
	db := NewPages()
	db.Open(dbPath)
	defer db.Close()

	AppendData(db.Pages[dbPath+"0"], "000000000000 ")

	getedData := GetValue(db.Pages[dbPath+"0"], 30, 20)

	fmt.Println("data:", getedData)
}

// Pages are map of file names that store data
type Pages struct {
	Pages map[string]*os.File
}

// NewPages creates List of files db
func NewPages() *Pages {
	return &Pages{
		Pages: make(map[string]*os.File, 1),
	}
}

// Opendb opnens all pages in root db
func (db *Pages) Open(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("readDir: ", err)
	}
	if len(files) == 0 {
		os.Create(path + "0")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		page, err := os.OpenFile(path+file.Name(), os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("os open file: ", err)
		}
		db.Pages[path+file.Name()] = page
		fmt.Println("file name is ", path+file.Name())
	}
}

// Close All pages
func (db *Pages) Close() {
	for _, Page := range db.Pages {
		Page.Close()
		fmt.Printf("%s closed\n", Page.Name())
	}
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

// NewPage creates new file db page
func NewPage(id int) {

	filename, _ := getLocation(id)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < 1000; i++ {
		// make spaces for indexes
		file.WriteString("               ") // lenght spaces 15
	}
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

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last + 1
}
