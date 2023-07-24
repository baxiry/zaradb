package dblite

import "os"

const IndexsFilePath = "primary.index"

var IndexsFile *os.File

func initIndexFile() {
	IndexsFile, _ = os.Open(IndexsFilePath)
	_ = IndexsFile
}

var pages = NewPages()

func initPages() {
	pages.Open(RootPath)
}

func init() {
	// check & init index map & firs page store
	initPages()
}
