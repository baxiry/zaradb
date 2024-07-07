package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func getCollections() string {
	table, result := "", `["`
	res, err := db.db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Println(err)
	}

	for res.Next() {
		res.Scan(&table)
		result += table + ", "
	}
	result = `{"collections": ` + result[:len(result)-2] + `"], "size": "123mb"}`

	return result
}

// deletes collection
func deleteCollection(query gjson.Result) string {
	// TODO return number of deleted objects
	return "not implemented yet"
}

// creates new collection
func createCollection(query gjson.Result) string {
	return "not implemented yet"
}

// Rename renames db.
func renameDB(query gjson.Result) error {
	return nil
}

// Remove remove db to .Trash dir
func removeDB(query gjson.Result) (err error) {
	return nil
}

// CreateDB create db. TODO return this directly
func createDB(query gjson.Result) (string, error) {

	return "not yet", nil
}

// DeleteDB deletes db. (free hard drive).
func deleteDB(query gjson.Result) string {
	return " not yet"
}
