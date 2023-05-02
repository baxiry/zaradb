package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const RootPath = "/Users/fedora/repo/dbs/"

// map of name files
type Pages struct {
	Pages map[string]*os.File
}

// creates List of files db
func NewPages() *Pages {
	return &Pages{
		Pages: make(map[string]*os.File, 1),
	}
}

// opnens all pages in root db
func (db *Pages) Open(path string) {
	fmt.Println("path is : ", path)

	// check primary.index file if exest
	_, err := os.Stat(path + IndexsFile)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.OpenFile(path+IndexsFile, os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("Error when create indes file", err)
			return
		}

	}

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("readDir: ", err)
	}
	if len(files) == 0 {
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
		fmt.Println("file name is ", path+file.Name())
	}

}

// closes All pages
func (db *Pages) Close() {
	for _, page := range db.Pages {
		page.Close()
		fmt.Printf("%s closed\n", page.Name())
	}
}

// creates new file db page
func (pages *Pages) NewPage(id int) {

	filename, _ := GetAt(id)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	sid := strconv.Itoa(id)

	pages.Pages[RootPath+sid] = file
	fmt.Printf("new page is created with %s name\n", RootPath+sid)

}
