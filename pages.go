package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const RootPath = "../dbs/"

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
	// check primary.index file if exest
	_, err := os.Stat(path + IndexsFile)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("indexes file not exist")
		// path/to/whatever does *not* exist

	}
	os.Exit(0)

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

// NewPage creates new file db page
func (pages *Pages) NewPage(id int) {

	filename, _ := GetLocation(id)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < 1000; i++ {
		// make spaces for indexes
		file.WriteString("               ") // lenght spaces 15
	}

	sid := strconv.Itoa(id)

	pages.Pages[RootPath+sid] = file
	fmt.Printf("new page is created with %s name\n", RootPath+sid)

}
