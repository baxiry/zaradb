package dblite

import (
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

func (cachedIndexs *CachedIndexs) GetIndex(id int) (pageName string, index [2]int64) {
	return strconv.Itoa(int(id) / 1000), cachedIndexs.indexs[id]
}

// initialize cache of indexs
func NewCachedIndexs() *CachedIndexs {

	cachedIndexs := &CachedIndexs{
		indexs: make([][2]int64, 0),
	}

	indxBuffer := make([]byte, IndexChnucLen)

	for {

		n, err := pages.Pages[indexFilePath].Read(indxBuffer)
		if err != nil && err != io.EOF {

			fmt.Println("ERROR! wher os.Read primary.Index file", err)
			fmt.Println("index file is ", pages.Pages[indexFilePath])
			os.Exit(1)
		}
		if err == io.EOF {
			break
		}

		slicIndexe := strings.Split(string(indxBuffer[:n]), " ")

		fmt.Println("length of slicIndexe: ", len(slicIndexe))

		at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
		size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

		cachedIndexs.indexs = append(cachedIndexs.indexs, [2]int64{at, size})
	}

	return cachedIndexs
}

// LastIndex return last index in table
func lastIndex(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		// TODO
		return 0 // panic("ERROR! no primary.index file ")
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

	IndexsCache.indexs = append(IndexsCache.indexs, [2]int64{int64(at), int64(dataSize)})
}

// deletes index from primary.index file
func DeleteIndex(id int, indxfile *os.File) { //
	at := int64(id * 20)
	indxfile.WriteAt([]byte("                    "), at)

	// TODO delete index from indexCache
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
