package db

import (
	"os"
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
		if f.IsDir() || !strings.HasSuffix(f.Name(), "primary index") {
			continue
		}
		nfiles += f.Name()[:len(f.Name())-2] + "\n"
	}
	// TODO count size of collections

	return nfiles
}

// deletes collection
func DeleteCollection(query string) string {
	collection := gjson.Get(query, "collection").String()
	// TODO return number of deleted objects
	return collection + " is deleted"
}

// creates new collection
func CreateCollection(query string) string {

	collection := gjson.Get(query, "collection").String()

	return "collecteon " + collection + " is created"
}
