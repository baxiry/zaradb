package main

import "os"

// CreateDB create db TODO return this directly
func CreateDB(dbName string) (err error) {
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {
	return os.MkdirAll(rootPath+dbName+"/.Trash/", 0755)
	// }
	// return err
}

// DeleteDB delete db (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " db deleted!"
}

// Remove remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

// Rename rename db.
func RenameDB(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}
