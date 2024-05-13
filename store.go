package store

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/wal"
)

var slash = Slash()

type Database struct {
	path        string
	Collections map[string]*Collection
}

type Collection struct {
	name      string
	log       *wal.Log
	indexs    []uint64
	lastIndex uint64
	id        uint64
}

// read index from data & build indexs
func (db *Database) buildIndexs() (indexs []uint64) {
	indexs = make([]uint64, 1)

	indexs = append(indexs, 1)

	for key, coll := range db.Collections {

		fmt.Println(" ++++ collection: ", key)

		for _, k := range coll.indexs[1:] {

			d, err := coll.getData(k)
			if err != nil {
				println("key is :", k, err.Error())
				return indexs
			}
			fmt.Printf("key %d data: %s \n", k, d)
		}
	}

	return indexs
}

func (coll *Collection) markDelete(id uint64) error {
	coll.lastIndex++
	err := coll.log.Write(coll.lastIndex, []byte("0"+strings.Repeat(" ", 19)))
	if err != nil {
		coll.lastIndex--
		return err
	}
	coll.indexs[id] = 0
	return err
}

// insert appends data
func (coll *Collection) insert(data string) error {
	coll.lastIndex++
	coll.id++
	id := fmt.Sprint(coll.id)
	err := coll.log.Write(coll.lastIndex, []byte(id+strings.Repeat(" ", 20-len(id))+data))
	if err != nil {
		//l, _ := coll.log.LastIndex()
		//println(err.Error(), coll.lastIndex, l)
		coll.lastIndex--
		coll.id--
		return err
	}
	coll.indexs = append(coll.indexs, coll.id)
	return nil
}

// getIndex get parse and return index or err
func (coll *Collection) getIndex(id uint64) (uint64, error) {
	bdata, err := coll.log.Read(id)
	if err != nil {
		return 0, err
	}
	index := strings.Split(string(bdata)[:20], " ")[0]
	i, err := strconv.Atoi(index)
	if err != nil {
		return 0, err
	}

	return uint64(i), nil
}

// getData reads data from wall file
func (coll *Collection) getData(id uint64) (string, error) {
	//println(len(coll.indexs), "id: ", id)
	//coll.indexs[id]
	bdata, err := coll.log.Read(id)
	if err != nil {
		return "", err
	}
	return string(bdata)[20:], nil
}

// NewEngine open exests path or creates new if not exists.
// this func create test as default collection
func NewDatabase(path string) *Database {
	db := &Database{
		Collections: make(map[string]*Collection, 0),
		path:        path,
	}

	dirs, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		// test is a default collection
		err = os.MkdirAll(db.path+"/test", 0766)
		if err != nil {
			panic(err)
		}

		db.NewCollection("test")
	} else {
		//		return nil
	}

	for _, p := range dirs {
		if p.IsDir() {
			db.NewCollection(p.Name())
		}
	}
	return db
}

func (db *Database) NewCollection(name string) error {

	coll := &Collection{name: name}
	log, err := wal.Open(db.path+slash+name, nil)
	if err != nil {
		return err
	}
	coll.log = log

	lastIndex, err := log.LastIndex()
	if err != nil {
		return err
	}

	coll.lastIndex = lastIndex
	coll.id = lastIndex
	if lastIndex == 0 {
		lastIndex++
	}
	coll.indexs = make([]uint64, lastIndex)

	db.Collections[name] = coll
	return nil
}

func (db *Database) Close() error {
	// close all collections
	for _, c := range db.Collections {
		err := c.log.Close()
		if err != nil {
			return err
		}
	}
	clear(db.Collections)
	return nil
}
