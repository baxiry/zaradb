package dblite

import (
	"errors"
	"os"
)

func initIndexsFile() {
	// check if primary.index is exist
	indexFilePath := db.Name + "testpi" //db.Index + pix
	_, err := os.Stat(indexFilePath)
	if errors.Is(err, os.ErrNotExist) {
		IndexsFile, err := os.OpenFile(indexFilePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			eLog.Println("when create indexFile.", err)
			return
		}
		//db.Pages[indexFilePath] = IndexsFile
		IndexsFile.Close()
	}
	//iLog.Println("indexFilePath is ", indexFilePath)
}

func initIndex() {
	indexFilePath := db.Name + "testpi" // db.Index + pix
	//collect = NewIndex("test")
	collect = InitIndex()
	collect.primaryIndex = lastIndex(indexFilePath)
	//db.colletions["testpi"].primaryIndex = lastIndex(indexFilePath)
}

//end
