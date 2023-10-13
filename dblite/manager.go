package dblite

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

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

	nfiles := "\n"
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
	collectName := gjson.Get(query, "collection").String()

	if collectName == "" {
		return "please type a collection you want to delete"
	}

	// remove all
	i := 0
	for range db.Pages {
		file := collectName + strconv.Itoa(i)
		i++
		for path := range db.Pages {
			if path == db.Name+file {
				_ = os.Remove(path)
				delete(db.Pages, path)
			}
		}
	}

	_ = os.Remove(db.Name + collectName + pIndex)
	delete(db.Pages, db.Name+collectName+pIndex)

	return collectName + " is deleted"
}

// creates new collection
func CreateCollection(query string) string {
	collectName := gjson.Get(query, "collection").String()

	if collectName == "" {
		return "please shose a collectin name"
	}

	// create index of collection & first page.

	firstPage, err := os.OpenFile(db.Name+collectName+fmt.Sprint(0), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}
	db.Pages[db.Name+collectName+fmt.Sprint(0)] = firstPage

	indxPage, err := os.OpenFile(db.Name+collectName+pIndex, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}
	db.Pages[db.Name+collectName+pIndex] = indxPage

	return "collecteon " + collectName + " is created"
}
