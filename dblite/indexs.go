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

// [[0,3],[3,8]]
type CachedIndexs struct {
	indexs [][2]int64
}

var IndexsCache *CachedIndexs

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
	db.PrimaryIndex = lastIndex(indexFilePath)
	IndexsCache = NewCachedIndexs()

	println("initialize Cached indexs, length is  ", len(IndexsCache.indexs))
}

func (cachedIndexs *CachedIndexs) GetIndex(id int) (pageName string, index [2]int64) {
	return strconv.Itoa(int(id) / 1000), cachedIndexs.indexs[id]
}

// initialize cache of indexs
func NewCachedIndexs() *CachedIndexs {
	path := db.Name + db.Collection + pi

	cachedIndexs := &CachedIndexs{
		indexs: make([][2]int64, 0),
	}

	indxBuffer := make([]byte, IndexChnucLen)

	for {
		//iLog.Println("indexFilePath: ", path)
		// iLog.Println("len of pages : ", len(db.Pages))

		n, err := db.Pages[path].Read(indxBuffer)
		if err != nil && err != io.EOF {
			eLog.Printf("ERROR! wher os.Read %s file %v", path, err)
			iLog.Println("index file is ", db.Pages[path])
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}

		slicIndexe := strings.Split(string(indxBuffer[:n]), " ")

		at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
		size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

		cachedIndexs.indexs = append(cachedIndexs.indexs, [2]int64{at, size})
	}
	iLog.Println("primary indexs length : ", len(cachedIndexs.indexs))

	At = cachedIndexs.lastAt()

	return cachedIndexs
}

// get last data location
func (cachedIndexs *CachedIndexs) lastAt() int {
	/*
		info, _ := os.Stat(db.Name + db.Collections + "0")
		fmt.Println("At is : ", info.Size())

	*/
	if len(cachedIndexs.indexs) > 0 {
		at := int(cachedIndexs.indexs[len(cachedIndexs.indexs)-1][0] + cachedIndexs.indexs[len(cachedIndexs.indexs)-1][1])
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
func AppendIndex(indexFile *os.File, at int, dataSize int) {

	strInt := fmt.Sprint(at) + " " + fmt.Sprint(dataSize)

	numSpaces := IndexChnucLen - len(strInt)
	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	//indexFile.WriteString(strInt)
	_, err := indexFile.WriteAt([]byte(strInt), db.PrimaryIndex*20)
	if err != nil {
		fmt.Println("err when UpdateIndex, store.go line 127", err)
	}

	IndexsCache.indexs = append(IndexsCache.indexs, [2]int64{int64(at), int64(dataSize)})
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
	IndexsCache.indexs[id] = [2]int64{dataAt, dataSize}
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
