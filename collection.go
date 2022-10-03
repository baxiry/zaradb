package main

import "os"

// CreateDB create db TODO return this directly
func CreateCl(coName string) (err error) {
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {
	return os.MkdirAll(rootPath+coName+"/.Trash/", 0755)
	// }
	// return err
}

// Delete delete db (free hard drive).
func DeleteCl(coName string) string {
	return coName + " collection deleted!"
}

// Remove remove db to .Trash dir
func RemoveCl(coName string) (err error) {
	return co.Rename(coName, ".Trash/"+coName)
}

// Rename rename db.
func RenameCl(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}
