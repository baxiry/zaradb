package main

import (
	"os"
)

func getCollections(dbName string) (collections []string) {
	return []string{"not implemented yet"}
}

// shows collections in corrent database
func showCollections(dbName string) string {
	return "not implemented yet"
}

// deletes collection
func deleteCollection(query string) string {
	// TODO return number of deleted objects
	return "not implemented yet"
}

// creates new collection
func createCollection(query string) string {
	return "not implemented yet"
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
