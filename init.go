package dblite

import (
	"errors"
	"fmt"
	"os"
)

var indexFilePath = RootPath + "primary.index"

var PrimaryIndex = lastIndex(indexFilePath)

func initIndexsFile() {
	// check if primary.index is exist
	_, err := os.Stat(indexFilePath)
	if errors.Is(err, os.ErrNotExist) {
		IndexsFile, err := os.OpenFile(indexFilePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		IndexsFile.Close()
	}
}

func initIndex() {
	lindx := lastIndex(indexFilePath)
	fmt.Println("last index is ", lindx)
}

func init() {
	// check & init index map & firs page store
	initIndex()

	initIndexsFile()
}
