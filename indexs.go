package dblite

import (
	"fmt"
	"io"
	"log"
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

func (cachedIndexs *CachedIndexs) GetIndex(id int) (pageName string, index [2]int64) {
	return strconv.Itoa(int(id) / 1000), cachedIndexs.indexs[id]
}

// initialize cache of indexs
func NewCachedIndexs() *CachedIndexs {
	cachedIndexs := &CachedIndexs{}
	indexs := make([][2]int64, 0)

	// TODO reuse
	ixBuffer := make([]byte, IndexChnucLen)

	for {

		n, err := pages.Pages[indexFilePath].Read(ixBuffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}

		slicIndexe := strings.Split(string(ixBuffer[:n]), " ")

		at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
		size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

		indexs = append(indexs, [2]int64{at, size})
	}
	cachedIndexs.indexs = indexs

	println("initialize Cached indexs ")

	return cachedIndexs
}

// LastIndex return last index in table
func lastIndex(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size() / 20
}

// append new index in primary.index file
func NewIndex(at int, dataSize int, indexFile *os.File) {

	strInt := fmt.Sprint(at) + " " + fmt.Sprint(dataSize)

	numSpaces := IndexChnucLen - len(strInt)
	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	indexFile.WriteString(strInt)

	// TODO add new index to chachedIndexs

	indexsCache.indexs = append(indexsCache.indexs, [2]int64{int64(at), int64(dataSize)})
}

// update index val in primary.index file
func UpdateIndex(id int, indexData, size int64, indexFile *os.File) {

	at := int64(id) * 20

	strIndex := fmt.Sprint(indexData) + " " + fmt.Sprint(size)
	for i := len(strIndex); i < 20; i++ {
		strIndex += " "
	}

	_, err := indexFile.WriteAt([]byte(strIndex), at)
	if err != nil {
		panic(err)
	}

}

// deletes index from primary.index file
func DeleteIndex(id int, indxfile *os.File) { //
	at := int64(id * 20)
	indxfile.WriteAt([]byte("                    "), at)
}

// get pageName Data Location  & data size from primary.indexes file
func GetIndex(id int, indexFile *os.File) (pageName string, at, size int64) {

	pageName = strconv.Itoa(int(id) / 1000)
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
