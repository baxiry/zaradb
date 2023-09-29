package dblite

import (
	"errors"
	"os"
)

func initIndexsFile() {
	// check if primary.index is exist
	indexFilePath := db.Name + db.Collection + pix
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

	iLog.Println("indexFilePath is ", indexFilePath)

}

func initIndex() {
	indexFilePath := db.Name + db.Collection + pix
	//collect = NewCollection("test")
	collect = InitCollection()
	collect.primaryIndex = lastIndex(indexFilePath)
}

/*
func init() {

	log.SetFlags(0) // Remove the default flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}
*/
