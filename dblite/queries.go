package dblite

import (
	"os"

	"github.com/tidwall/gjson"
)

// ok
func HandleQueries(query string) string {
	switch gjson.Get(query, "action").String() {

	case "findOne":
		return findOne(query)

	case "findMany":
		return findMany(query)

	case "findById":
		return findById(query)

	case "insert":
		return Insert(query)

	case "update":
		return Update(query)

	case "delete":
		return DeleteById(query)

	// manage database
	case "create_collection":
		return CreateCollection(query)

	case "delete_collection":
		return DeleteCollection(query)

	case "show_collection":
		return showCollections(db.path)

	default:
		return "unknowen action"
	}
}

// extract filter
func findByField(json, field string) string {
	return gjson.Get(json, field).String()
}

// Rename renames db.
func RenameDB(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Remove remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

// CreateDB create db. TODO return this directly
func CreateDB(dbName string) (string, error) {

	return dbName + "is created", nil
}

// DeleteDB deletes db. (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " is deleted!"
}
