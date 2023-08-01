package dblite

import (
	"errors"
	"fmt"
	"os"
)

type Enginge struct{}

var pages = NewPages()

func (Enginge) Run() {
	println("Engine Runing...")

	// check & init index map & firs page store

	initIndexsFile()

	initPages()

	initIndex()

}

func (Enginge) Stop() {
	println("Enginge Closing...")
	pages.Close()
}

func NewEngine() *Enginge {
	return &Enginge{}
}

var indexFilePath = RootPath + "primary.index"

var PrimaryIndex = lastIndex(indexFilePath)

var IndexsCache *CachedIndexs

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
	_ = lastIndex(indexFilePath)
	IndexsCache = NewCachedIndexs()

	println("initialize Cached indexs length is  ", len(IndexsCache.indexs))
}

func initPages() {
	file, err := os.OpenFile(RootPath+"0", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {

		fmt.Println("initIndexsFile, OpenFile")
		panic(err)
	}
	file.Close()

	pages.Open(RootPath)

}
