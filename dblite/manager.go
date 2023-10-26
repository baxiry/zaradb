package dblite

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func getCollections(dbName string) (collections []string) {
	f, err := os.Open(dbName)
	if err != nil {
		eLog.Println(err)
		return nil
	}

	files, err := f.Readdir(0)
	if err != nil {
		eLog.Println(err)
		return nil
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), pIndex) {

			continue
		}
		collections = append(collections, f.Name())
	}

	return collections
}

// shows collections in corrent database
func showCollections(dbName string) string {
	f, err := os.Open(dbName)
	if err != nil {
		return err.Error()
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err.Error()
	}

	nfiles := ""
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), pIndex) {
			continue
		}
		nfiles += f.Name()[:len(f.Name())-2] + "\n"
	}
	// TODO count size of collections

	return nfiles
}

// deletes collection
func DeleteCollection(query string) string {
	collection = gjson.Get(query, "collection").String()

	if collection == "" {
		return "please type a collection you want to delete"
	}

	// remove all
	i := 0
	for range db.Pages {
		file := collection + strconv.Itoa(i)
		i++
		for path := range db.Pages {
			if path == db.Name+file {
				_ = os.Remove(path)
				delete(db.Pages, path)
			}
		}
	}
	_ = os.Remove(db.Name + collection + pIndex)
	delete(db.Pages, db.Name+collection+pIndex)

	// TODO return number of deleted objects
	return collection + " is deleted"
}

// creates new collection
func CreateCollection(query string) string {

	collection = gjson.Get(query, "collection").String()
	action := gjson.Get(query, "action").String()

	if collection == "" {
		if action != "" {
			return "please shose a collectin name"
		}
		collection = query
	}

	// create index of collection & first page.

	firstPage, err := os.OpenFile(db.Name+collection+fmt.Sprint(0), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}

	ipath := db.Name + collection + pIndex

	db.Pages[db.Name+collection+fmt.Sprint(0)] = firstPage

	indxPage, err := os.OpenFile(ipath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}
	db.Pages[ipath] = indxPage

	// init new index in indexs

	Indexs[ipath] = &Index{at: 0, indexCache: [][2]int64{}, primaryIndex: 0}
	//	iLog.Println(Indexs)

	return "collecteon " + collection + " is created"
}
