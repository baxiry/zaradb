package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// delete
// TODO database is locked
func deleteOne(query string) string {

	coll := gjson.Get(query, "collection").String()
	filter := gjson.Get(query, "filter").String()

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	rowid := "0"
	record := ""

	for rows.Next() {
		record = ""
		rowid = ""
		err := rows.Scan(&rowid, &record)
		if err != nil {
			fmt.Println(err.Error())
		}
		if match(filter, record) {
			break
		}
	}

	rows.Close()

	res, err := db.db.Exec(`delete from ` + coll + ` where rowid = ` + rowid) // fast
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	println("res: ", res)

	return `{"result":"` + rowid + ` removed"}`

	//return `{"result":"no match to removed"}`
}

// Finds first obj match creteria.
func findById(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	id := gjson.Get(query, "_id").String()

	stmt := `select record from ` + coll + ` where rowid = ` + id

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	record := ""
	for rows.Next() {
		err := rows.Scan(&record)
		if err != nil {
			return err.Error()
		}
		return record
	}

	return `{"result":0}`
}

// Finds first obj match creteria.
func findOne(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	filter := gjson.Get(query, "filter").String()

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	record := ""
	for rows.Next() {
		err := rows.Scan(&record)
		if err != nil {
			return err.Error()
		}
		if match(filter, record) {
			return record
		}
	}

	return `{"result":0}`
}

// Find finds any obs match creteria.
func findMany(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	filter := gjson.Get(query, "filter").String()

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	records := "["
	record := ""

	for rows.Next() {
		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return err.Error()
		}
		if match(filter, record) {
			records += record + `,`
		}
	}
	return records[:len(records)-1] + "]"
}

// Insert
func insert(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").String()

	err := db.insert(coll, data)
	if err != nil {
		return err.Error()
	}

	return "inser done"
}

// Insert
func insertMany(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").String()

	err := db.insert(coll, data)
	if err != nil {
		return err.Error()
	}

	return "inser done"
}

// delete
func deleteMany(query string) string {

	coll := gjson.Get(query, "collection").String()
	filter := gjson.Get(query, "filter").String()

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	records := "["
	record := ""

	for rows.Next() {
		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return err.Error()
		}
		if match(filter, record) {
			records += record + `,`
		}
	}
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
