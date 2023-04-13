package main

import (
	"fmt"
	"os"
	"strings"
)

// simplest query language
func queryLang() {
	query := arguments()
	fmt.Println("query is : ", query)
}

// Insert
func Insert(path, data string) (err error) {
	return
}

// Select reads data form docs
func Select(path string) (data string) {
	return data
}

// Update update document data
func Update(serial, data string) (err error) {
	return
}

// Delete removes document
func Delete(path string) (err error) {
	return
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

// helpers function

// check if path is exist
func IsExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// ListDir show all directories in path
func ListDir(path string) {
	dbs, err := os.ReadDir(RootPath + path)
	if err != nil {
		fmt.Println(err)
	}

	dirs := 0
	for _, dir := range dbs {
		if dir.IsDir() && string(dir.Name()[0]) != "." {
			dirs++
			print(dir.Name(), " ")
		}
	}
	if dirs > 0 {
		println()
		return
	}
	println(path, "is impty")
}

// PathExist check if path exists
func PathExist(subPath string) bool {
	_, err := os.Stat(RootPath + subPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
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

	err := os.MkdirAll(RootPath+dbName+"/.Trash/", 0755)
	if err != nil {
		return dbName, err
	}
	return dbName, nil
}

// DeleteDB deletes db. (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " db deleted!"
}
