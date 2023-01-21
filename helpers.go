package main

import (
	"fmt"
	"os"
)

// ListDir prints all directories
func ListDir(path string) {
	dbs, err := os.ReadDir(rootPath + path)
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

// PathExist check if path exists & return boolean
func PathExist(subPath string) bool {
	_, err := os.Stat(rootPath + subPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// LastIndex return last index in table
func LastIndex(path string) int {
	return 0
}

// Rename renames db.
func RenameDB(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Remove remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

// CreateDB create db TODO return this directly
func CreateDB(dbName string) (string, error) {
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {return err}

	err := os.MkdirAll(rootPath+dbName+"/.Trash/", 0755)
	if err != nil {
		return dbName, err
	}
	return dbName, nil
}

// DeleteDB delete db (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " db deleted!"
}
