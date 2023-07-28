package dblite

import (
	"fmt"
	"os"
	"strconv"
)

// Root database folder
var RootPath string = userDir() + "/repo/dbs/"

// Mock path
var MockPath string = userDir() + "/repo/dblite/mok/"

var indexsCache *CachedIndexs

// map of name files
type Pages struct {
	Pages map[string]*os.File
}

// NewPages constracts List of files pages
func NewPages() *Pages {
	return &Pages{
		Pages: make(map[string]*os.File, 2),
	}
}

// creates new page file and add it to Pages Map
func (pages *Pages) NewPage(id int) {
	// TODO

	filename, _, _ := GetIndex(id, pages.Pages[indexFilePath])

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	// do not close this file

	strId := strconv.Itoa(id)

	pages.Pages[RootPath+strId] = file
	// fmt.Printf("new page is created with %s name\n", RootPath+strId)

}

// opnens all pages in Root database folder
func (pages *Pages) Open(path string) {

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("readDir: ", err)
		return
	}

	if len(files) < 1 {
		//os.Create(path + "0")
		os.OpenFile(path+"primary.index", os.O_CREATE|os.O_RDWR, 0644)
		// return
	}
	if len(files) < 2 {
		//os.Create(path + "0")
		os.OpenFile(path+"0", os.O_CREATE|os.O_RDWR, 0644)
		// return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		page, err := os.OpenFile(path+file.Name(), os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("os open file: ", err)
		}
		pages.Pages[path+file.Name()] = page
	}
	fmt.Println("pages is ready")

	indexsCache = NewCachedIndexs()
	fmt.Println("icache is ready")

}

// closes All pages
func (pages *Pages) Close() {
	for _, page := range pages.Pages {
		page.Close()
		fmt.Printf("%s closed\n", page.Name())
	}
}
