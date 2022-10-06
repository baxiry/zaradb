package main

import "os"

// type Collection string

// CreateDB create db TODO return this directly
func CreateCl(cPath string) (colname string, err error) { // db and collection Path
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {
	err = os.MkdirAll(cPath+"/.Trash/", 0755)
	return colname, err
}

// Delete delete db (free hard drive).
func DeleteCl(cPath string) string {
	return cPath + " collection deleted!"
}

// Remove remove db to .Trash dir
func RemoveCl(cPath string) (err error) {
	return RenameCl(cPath, ".Trash/"+cPath)
}

// Rename rename db.
func RenameCl(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}
