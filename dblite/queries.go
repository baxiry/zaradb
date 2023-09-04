package dblite

import (
	"fmt"
	"os"
	"strings"

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

	default:
		return "unknowena ction"

	}

	//return result.String()
}

// extractQuery from stdin argument
func extractQuery(str string) (json string) {
	var start, end int32

	var i int32
	for i = 0; i < int32(len(str)); i++ {
		if str[i] == '{' {
			start = i
			break
		}
	}
	for i = int32(len(str)) - 1; i >= 0; i-- {
		if str[i] == '}' {
			end = i
			break
		}
	}
	return str[start : end+1]
}

// cli functions

const hints = `tap helpe to get help massage`

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return strings.Split(args[1], ".")
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
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {return err}

	err := os.MkdirAll(db.Name+"/.Trash/", 0755)
	if err != nil {
		return dbName, err
	}
	return dbName, nil
}

// DeleteDB deletes db. (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " is deleted!"
}
