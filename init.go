package dblite

import (
	"fmt"
)

var PrimaryIndex = lastIndex(indexFilePath)

var indexFilePath = RootPath + "primary.index"

/*
var IndexsFile *os.File

func initIndexsFile() {

		// check if primary.index is exist
		_, err := os.Stat(indexFilePath)
		if errors.Is(err, os.ErrNotExist) {
			IndexsFile, err = os.Open(indexFilePath)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("IndexsFile wher path ", IndexsFile)
	}
*/
func initIndex() {
	lindx := lastIndex(indexFilePath)
	fmt.Println("last index is ", lindx)
}

func init() {
	// check & init index map & firs page store
	initIndex()

	// initIndexsFile()
}
