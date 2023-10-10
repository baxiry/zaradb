package dblite

import (
	"os"

	"github.com/tidwall/gjson"
)

func HandleQueries(query string) string {
	action := gjson.Get(query, "action")

	switch action.String() {
	case "insert":
		return Insert(query)

	case "select":

		return SelectById(query)

	case "update":

		return Update(query) // db.Name = RootPath

	case "delete":
		DeleteById(query)
		return "action is Delete"

	case "create_collection":
		return CreateCollection(query)

	case "delete_collection":
		return DeleteCollection(query)

	default:
		return "unknowen action"

	}

	//return result.String()
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
