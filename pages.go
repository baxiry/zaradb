package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Root database folder
var RootPath string = userDir() + "/repo/dbs/"

// map of name files
type Pages struct {
	Pages map[string]*os.File
}

// NewPages constracts List of files db
func NewPages() *Pages {
	return &Pages{
		Pages: make(map[string]*os.File, 1),
	}
}

// creates new page file and add it to Pages Map
func (pages *Pages) NewPage(id int) {

	filename, _ := GetWhere(id)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	// do not close this file

	strId := strconv.Itoa(id)

	pages.Pages[RootPath+strId] = file
	fmt.Printf("new page is created with %s name\n", RootPath+strId)

}

// opnens all pages in Root database folder
func (db *Pages) Open(path string) {
	indexFile := path + IndexsFile

	// check if primary.index is exist
	_, err := os.Stat(indexFile)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.OpenFile(indexFile, os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("Error when create indes file", err)
			return
		}

	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("readDir: ", err)
	}
	if len(files) < 2 {
		os.Create(path + "0")
		return
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
		fmt.Println("file name is ", path+file.Name(), "is Open")
	}
}

// closes All pages
func (db *Pages) Close() {
	for _, page := range db.Pages {
		page.Close()
		fmt.Printf("%s closed\n", page.Name())
	}
}
