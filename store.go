package store

import (
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/wal"
)

func (coll *Collection) insert(data string) error {
	coll.lastIndex++
	coll.id++
	id := fmt.Sprint(coll.id)
	err := coll.log.Write(coll.lastIndex, []byte(id+strings.Repeat(" ", 20-len(id))+data))
	if err != nil {
		l, _ := coll.log.LastIndex()
		println(err.Error(), coll.lastIndex, l)
		coll.lastIndex--
		coll.id--
		return err
	}
	coll.indexs = append(coll.indexs, coll.id)
	return nil
}

func (coll *Collection) get(id uint64) (string, error) {
	bdata, err := coll.log.Read(id)
	if err != nil {
		return "", err
	}
	return string(bdata)[20:], nil
}

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

// NewEngine open exests path or creates new if not exists.
// this func create test as default collection
func NewDatabase(path string) *Database {
	db := &Database{
		Collections: make(map[string]*Collection, 0),
		path:        "../dbs/" + path,
	}

	dirs, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		// test is a default collection
		err = os.MkdirAll(db.path+"/test", 0766)
		if err != nil {
			panic(err)
		}

		db.newCollection("test")
	} else {
		//		return nil
	}

	for _, p := range dirs {
		if p.IsDir() {
			db.newCollection(p.Name())
		}
	}
	return db
}

func (db *Database) newCollection(name string) error {

	coll := &Collection{name: name}
	log, err := wal.Open(db.path+"/"+name, nil)
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
