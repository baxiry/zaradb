package dblite

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	// pix is primary index file
	pIndex = "pi"

	// buffer size of len
	IndexChnucLen = 20
)

type Index struct { // Index
	// at : data locations store in file
	at int64
	// current primaryIndex value
	primaryIndex int64
	// indexes cache
	indexCache [][2]int64 // [[0,3],[3,8]]
}

// list of collections
var Indexs = make(map[string]*Index)

// global var! be careful
//var index = Index{}

// initialize cache of indexs
func InitIndex() map[string]*Index {
	indexs := make(map[string]*Index)

	// iLog.Println("indexFilePath: ", path)

	indxBuffer := make([]byte, IndexChnucLen)

	// get all collection in this database first
	//iFile := db.Name + collection + pIndex
	for _, indexfile := range getCollections(db.Name) {

		index := &Index{
			at:           0,
			primaryIndex: 0,
			indexCache:   make([][2]int64, 0),
		}

		for {
			n, err := db.Pages[db.Name+indexfile].Read(indxBuffer)
			if err != nil && err != io.EOF {
				eLog.Println("file is : ", err)
				os.Exit(1)
			}
			if err == io.EOF {
				break
			}
			if n%20 != 0 {
				eLog.Println("why n is :", n)
			}

			slicIndexe := strings.Split(string(indxBuffer[:n]), " ")

			fmt.Printf("slicIndexe: .%s.\n", string(indxBuffer[:n]))
			fmt.Println("path", db.Name+indexfile)
			// TODO check bug here
			if len(slicIndexe) == 1 {
				eLog.Println("len slicIndexe is just 1", slicIndexe[0])
				continue
			}
			fmt.Println()

			at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
			size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

			index.indexCache = append(index.indexCache, [2]int64{at, size})
		}

		iLog.Println("indexs length : ", len(indexs))
		fmt.Println()

		indexs[indexfile] = index

		lst, primary := lasts(db.Name + indexfile)
		indexs[indexfile].at = lst               // check here
		indexs[indexfile].primaryIndex = primary // check here
	}

	for k, v := range indexs {
		fmt.Printf("at in %s is %d\n", k, v.at)
		fmt.Printf("pi in %s is %d\n", k, v.primaryIndex)
	}

	return indexs
}

// GetIndex
func (c *Index) GetIndex(id int) (pageName string, index [2]int64) {
	return strconv.Itoa(int(id) / 1000), c.indexCache[id]
}

// get last data location
func lasts(path string) (int64, int64) {
	info, err := os.Stat(path)
	if err != nil {
		// TODO
		eLog.Printf(path, err)
		return 0, 0 // panic("ERROR! no primary.index file ")
	}

	size := info.Size()

	at := size - 20
	buf := make([]byte, 20)
	f, _ := os.OpenFile(path, os.O_RDONLY, 0644)
	f.ReadAt(buf, at)

	slc := strings.Split(string(buf), " ")
	lastat, _ := strconv.ParseInt(slc[0], 10, 64)

	lastPrimaryIndex := size / 20

	return lastat, lastPrimaryIndex
}

/*
// get last data location
func (c *Index) lastAt() int64 {
	if len(c.indexCache) > 0 {

		at := c.indexCache[len(c.indexCache)-1][0] + c.indexCache[len(c.indexCache)-1][1]
		println("last at is ", at)
		return at
	}
	return 0
}
*/

// LastIndex return last index in table
func lastIndex(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		// TODO
		eLog.Println("pi is not exists ")
		return 0 // panic("ERROR! no primary.index file ")
	}

	// for file := range Indexs {fmt.Println(file)}
	//iLog.Println("last index from Indexs Cache:  ", len(indexs[collection+pIndex].indexCache))
	//iLog.Println("last index from Index file Size", info.Size()/20)

	return info.Size() / 20
}

// append new index in pi file
func AppendIndex(indexFile *os.File, at int64, dataSize int) {

	strInt := fmt.Sprint(at) + " " + fmt.Sprint(dataSize)

	numSpaces := IndexChnucLen - len(strInt)
	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	iLog.Println("Collection in AppendIndex is : ", collection+pIndex)

	_, err := indexFile.WriteAt([]byte(strInt), Indexs[collection+pIndex].primaryIndex*20) // indexfile.Name()
	if err != nil {
		fmt.Println("err when UpdateIndex, store.go line 127", err)
	}

	Indexs[collection+pIndex].indexCache = append(Indexs[collection+pIndex].indexCache, [2]int64{at, int64(dataSize)})
	// TODO use assgined via index insteade append here e.g indexs[coll].indexs[id] = [2]int64{at, dataSize}
}

// update index val in primary.index file & cache index file
func UpdateIndex(indexFile *os.File, id int, dataAt, dataSize int64) {

	at := int64(id * 20) // shnuck

	strIndex := fmt.Sprint(dataAt) + " " + fmt.Sprint(dataSize) + " "

	_, err := indexFile.WriteAt([]byte(strIndex), at)
	if err != nil {
		fmt.Println("id & at is ", id, at)
		eLog.Println("err when UpdateIndex", err)
	}

	Indexs[collection].indexCache[id] = [2]int64{dataAt, dataSize}
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
