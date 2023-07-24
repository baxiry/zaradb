package dblite

import (
	"errors"
	"fmt"
	"os"
)

const IndexsFilePath = "primary.index"

var iFilePath = RootPath + IndexsFilePath

var pages = NewPages()

var IndexsFile *os.File

func initIndexsFile() {

	// check if primary.index is exist
	_, err := os.Stat(iFilePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err := os.OpenFile(iFilePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("Error when create indes file", err)
			return
		}

	}

}

func initIndex() {
	lastIndex(iFilePath)
	fmt.Println("index file is ", iFilePath)
}

func initPages() {
	pages.Open(RootPath)
}

func init() {
	// check & init index map & firs page store
	initPages()
	initIndex()
}
