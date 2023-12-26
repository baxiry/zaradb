package dblite

import (
	"os"

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
		nfiles += f.Name()[:len(f.Name())-2] + "\n"
	}
	// TODO count size of collections

	return nfiles
}

// deletes collection
func deleteCollection(query string) string {
	collection := gjson.Get(query, "collection").String()

	if collection == "" {
		return "please type a collection you want to delete"
	}

	// remove all
	for range db.pages {
	}

	// TODO return number of deleted objects
	return collection + " is deleted"
}

// creates new collection
func createCollection(query string) string {

	collection := gjson.Get(query, "collection").String()
	action := gjson.Get(query, "action").String()

	if collection == "" {
		if action != "" {
			return "please shose a collectin name"
		}
		collection = query
	}

	// create index of collection & first page.

	return "collecteon " + collection + " is created"
}

// Rename renames db.
func renameDB(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Remove remove db to .Trash dir
func removeDB(dbName string) (err error) {
	return renameDB(dbName, ".Trash/"+dbName)
}

// CreateDB create db. TODO return this directly
func createDB(dbName string) (string, error) {

	return dbName + "is created", nil
}

// DeleteDB deletes db. (free hard drive).
func deleteDB(dbName string) string {
	return dbName + " is deleted!"
}
