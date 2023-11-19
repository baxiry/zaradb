package dblite

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	// pIndex is primary index file sufix
	pIndex = "pi"

	// size of index buffer
	IndexChunkLen = 20
)

// Index stor at primaryIndex & indexCache
type Index struct { // Index
	// at : data locations store in file
	at int64
	// current primaryIndex value
	primaryIndex int64
	// indexes cache
	indexCache [][2]int64 // [[0,3],[3,8]]
}

// list of Index type
var Indexs = make(map[string]*Index)

// initialize cache of indexs
func InitIndex() map[string]*Index {
	indexs := make(map[string]*Index)

	indxBuffer := make([]byte, IndexChunkLen)

	// get all collection in this database first
	for _, indexfile := range getCollections(db.Name) {
		indexfile = db.Name + indexfile
		index := &Index{
			at:           0,
			primaryIndex: 0,
			indexCache:   make([][2]int64, 0),
		}

		for {
			n, err := db.Pages[indexfile].Read(indxBuffer)
			if err != nil && err != io.EOF {
				eLog.Printf("file: %s. %s\n", indexfile, err)
				os.Exit(1)
			}
			if err == io.EOF {
				break
			}

			slicIndexe := strings.Split(string(indxBuffer[:n]), " ")

			// TODO check bug here
			if len(slicIndexe) == 1 {
				eLog.Println("len slicIndexe is just 1", slicIndexe[0])
				continue
			}

			at, _ := strconv.ParseInt(slicIndexe[0], 10, 64)
			size, _ := strconv.ParseInt(slicIndexe[1], 10, 64)

			index.indexCache = append(index.indexCache, [2]int64{at, size})
		}

		indexs[indexfile] = index

		lst, primary := lasts(indexfile)
		indexs[indexfile].at = lst               // check here
		indexs[indexfile].primaryIndex = primary // check here
	}

	return indexs
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
	s, _ := strconv.ParseInt(slc[1], 10, 64)

	lastPrimaryIndex := size / 20

	return lastat + s, lastPrimaryIndex
}

// LastIndex return last index in table
func lastIndex(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		// TODO
		eLog.Println("pi is not exists ")
		return 0 // panic("ERROR! no primary.index file ")
	}

	return info.Size() / 20
}

// append new index in pi file
func AppendIndex(indexFile *os.File, at int64, dataSize int) {

	strInt := fmt.Sprint(at) + " " + fmt.Sprint(dataSize)

	numSpaces := IndexChunkLen - len(strInt)
	for i := 0; i < numSpaces; i++ {
		strInt += " "
	}

	ipath := db.Name + collection + pIndex

	_, err := indexFile.WriteAt([]byte(strInt), Indexs[ipath].primaryIndex*20) // indexfile.Name()
	if err != nil {
		eLog.Println("UpdateIndex", err)
	}

	Indexs[ipath].indexCache = append(Indexs[ipath].indexCache, [2]int64{at, int64(dataSize)})
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

	Indexs[db.Name+collection+pIndex].indexCache[id] = [2]int64{dataAt, dataSize}
}

// get pageName Data Location  & data size from primary.indexes file
func GetIndex(indexFile *os.File, id int) (pageName string, at, size int64) {

	pageName = strconv.Itoa(id / int(MaxObjects))
	bData := make([]byte, IndexChunkLen)
	_, err := indexFile.ReadAt(bData, int64(id*IndexChunkLen))
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

	// TODO SEARCH for how can i controll size of file without use string space ?
	indxfile.WriteAt([]byte("                    "), at)

	// TODO Delete index from indexCache or use a Bloom filter to avoid unnecessary file reading.
}

//end
