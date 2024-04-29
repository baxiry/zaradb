package main

import (
	"github.com/tidwall/gjson"
)

// Finds first obj match creteria.
func findOne(query string) (res string) {
	return "not implemented yet"
}

// Find finds any obs match creteria.
func findMany(query string) (res string) {

	return "not implemented yet"
}

// findById reads data form docs
func findById(query string) string {
	return "not implemented yet"
}

// Insert
func insert(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").String()

	err := db.Insert(coll, data)
	if err != nil {
		println("network problem") // network problem
		return err.Error()
	}

	return "not implemented yet"
}

// delete
func deleteOne(query string) string {
	return "not implemented yet"
}

// delete
func deleteMany(query string) string {

	return "not implemented yet"

}

// delete by id
func deleteById(query string) string {
	return "not implemented yet"
}

// Update update document data
func updateById(query string) (result string) {
	return "not implemented yet"
}

// TODO updateOne one update document data
func updateOne(query string) (result string) {
	return "not implemented yet"
}

// TODO updateMany update document data
func updateMany(query string) (result string) {
	return "not implemented yet"
}

// end
