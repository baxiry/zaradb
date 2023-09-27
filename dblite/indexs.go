package dblite

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// buffer size of len
const IndexChnucLen = 20

var collect *Collection

func initIndexsFile() {
	// check if primary.index is exist
	indexFilePath := db.Name + db.Collection + pi
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
	indexFilePath := db.Name + db.Collection + pi
	collect = InitCollection()
	collect.primaryIndex = lastIndex(indexFilePath)
}

func (c *Collection) GetIndex(id int) (pageName string, index [2]int64) {
	return strconv.Itoa(int(id) / 1000), c.cachedIndexs[id]
}

// initialize cache of indexs
func InitCollection() *Collection {
	path := db.Name + db.Collection + pi
	iLog.Println("indexFilePath: ", path)
	iLog.Println("len of pages : ", len(db.Pages))

	c := &Collection{
		cachedIndexs: make([][2]int64, 0),
	}

	indxBuffer := make([]byte, IndexChnucLen)

	for {
		n, err := db.Pages[path].Read(indxBuffer)
		if err != nil && err != io.EOF {
			eLog.Printf("ERROR! %s wher os.Read  file %v\n", err, path)
			iLog.Println("file is : ", db.Pages[path])
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}

		slicIndexe := strings.Split(string(indxBuffer[:n]), " ")

		at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
		size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

		c.cachedIndexs = append(c.cachedIndexs, [2]int64{at, size})
	}
	iLog.Println("primary indexs length : ", len(c.cachedIndexs))

	c.at = c.lastAt()

	return c
}

// get last data location
func (c *Collection) lastAt() int64 {
	if len(c.cachedIndexs) > 0 {
		at := c.cachedIndexs[len(c.cachedIndexs)-1][0] + c.cachedIndexs[len(c.cachedIndexs)-1][1]
		println("At is ", at)
		return at
	}
	return 0
}

// LastIndex return last index in table
func lastIndex(path string) int64 {
	iLog.Println("path in last index func is ", path)
	info, err := os.Stat(path)
	if err != nil {
		// TODO
		eLog.Println("pi is not exists ")
		return 0 // panic("ERROR! no primary.index file ")
	}

	iLog.Println("last index is", info.Size()/20)
	return info.Size() / 20
}

// append new index in pi file
func AppendIndex(indexFile *os.File, at int64, dataSize int) {

	strInt := fmt.Sprint(at) + " " + fmt.Sprint(dataSize)

	numSpaces := IndexChnucLen - len(strInt)
	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	//indexFile.WriteString(strInt)
	_, err := indexFile.WriteAt([]byte(strInt), collect.primaryIndex*20)
	if err != nil {
		fmt.Println("err when UpdateIndex, store.go line 127", err)
	}

	collect.cachedIndexs = append(collect.cachedIndexs, [2]int64{at, int64(dataSize)})
	// TODO use assgined insteade append here e.g IndexsCache.indexs[id] = [2]int64{int64(at), int64(dataSize)}
}

// update index val in primary.index file
func UpdateIndex(indexFile *os.File, id int, dataAt, dataSize int64) {

	at := int64(id * 20)

	strIndex := fmt.Sprint(dataAt) + " " + fmt.Sprint(dataSize) + " "
	//for i := len(strIndex); i < 20; i++ {	strIndex += " "}

	_, err := indexFile.WriteAt([]byte(strIndex), at)
	if err != nil {
		fmt.Println("id & at is ", id, at)
		fmt.Println("err when UpdateIndex, store.go line 127", err)

	}

	// TODO update index in indexsCache
	collect.cachedIndexs[id] = [2]int64{dataAt, dataSize}
}

// get pageName Data Location  & data size from primary.indexes file
func GetIndex(indexFile *os.File, id int) (pageName string, at, size int64) {

	pageName = strconv.Itoa(id / int(MaxObjects))
	bData := make([]byte, 20)
	_, err := indexFile.ReadAt(bData, int64(id*20))
	if err != nil {
		panic(err)
	}

	slc := strings.Split(string(bData), " ")
	iat, _ := strconv.Atoi(slc[0])

	isize, _ := strconv.Atoi(fmt.Sprint(slc[1]))
	return pageName, int64(iat), int64(isize)
}

// deletes index from primary.index file
func DeleteIndex(indxfile *os.File, id int) { //
	at := int64(id * 20)
	indxfile.WriteAt([]byte("                    "), at)
	// TODO delete index from indexCache
}

//end
