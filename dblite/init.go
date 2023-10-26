package dblite

import (
	"errors"
	"os"
)

func initIndexs() {
	// initialize Primaryindexs file if not exist cache
	indexFilePath := db.Name + testCollection + pIndex
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

	// initialize indexs cache
	Indexs = InitIndex()
}

//end
