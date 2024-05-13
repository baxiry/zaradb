package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// Find finds any obs match creteria.
func (db *DB) findMany(query string) (res string) {

	// TODO parse hol qury one time
	coll := gjson.Get(query, "collection").String()
	if coll == "" {
		return `{"error":"forgot collection name "}`
	}

	filter := gjson.Get(query, "filter").String()

	if filter == "" {
		filter = "{}"
	}

	skip := gjson.Get(query, "skip").Int()
	limit := gjson.Get(query, "limit").Int()
	if limit == 0 {
		limit = 100 // what is default setting ?
	}
	fmt.Println(limit, skip)

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	records := "["
	record := ""

	for rows.Next() {
		if limit == 0 {
			break
		}
		if skip != 0 {
			skip--
			continue
		}

		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return err.Error() // TODO standard errors
		}

		ok, err := match(filter, record)
		if err != nil {
			return err.Error()
		}

		if ok {
			// aggrigate here
			records += record + `,`
			limit--
		}
	}
	return records[:len(records)-1] + "]"
}

// Finds first obj match creteria.
func (db *DB) findOne(query string) (res string) {
	coll := gjson.Get(query, "collection").String()

	filter := gjson.Get(query, "filter").String()
	skip := gjson.Get(query, "skip").Int()

	// TODO are skyp useful here ?

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	record := ""
	for rows.Next() {
		if skip != 0 {
			skip--
			continue
		}

		err := rows.Scan(&record)
		if err != nil {
			return err.Error()
		}
		b, err := match(filter, record)
		if err != nil {
			return err.Error()
		}
		if b {
			return record
		}
	}

	return `{"result":0}`
}

// delete
func (db *DB) deleteOne(query string) string {

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

		b, err := match(filter, record)
		if err != nil {
			return err.Error()
		}
		if b {
			break
		}
	}
	// should close here
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
func (db *DB) findById(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	id := gjson.Get(query, "_id").String()

	stmt := `select record from ` + coll + ` where rowid = ` + id

	rows, err := db.db.Query(stmt)
	if err != nil {
		if id == "" {
			return `{Error:"Parameter not found: _id"}`
		}
		return fmt.Sprintf(`{Error:"%s"}`, err)
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

// Insert
func (db *DB) insertOne(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").String()

	err := db.insert(coll, data)
	if err != nil {
		return err.Error()
	}

	return "inser done"
}

// deleteMany
func (db *DB) deleteMany(query string) string {

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
		b, err := match(filter, record)
		if err != nil {
			return err.Error()
		}
		if b {
			records += record + `,`
		}
	}
	return "not implemented yet"

}

// delete by id
func (db *DB) deleteById(query string) string {
	id := gjson.Get(query, "_id").String()
	if id == "" {
		return `{"error": "there is no _id"}`
	}
	coll := gjson.Get(query, "collection").String()
	if coll == "" {
		return `{"error": "there is no collection"}`
	}

	sql := `delete from ` + coll + ` where rowid = ` + id
	fmt.Println(sql)

	_, err := db.db.Exec(sql)
	if err != nil {
		return `{"error": "internal error"}` // + err.Error()
	}
	return `{"aknowlge": "row ` + id + ` deleted"}`
}

// Update update document data
func (db *DB) updateById(query string) (result string) {
	data := db.findById(query)
	println("data: ", data)

	return "not implemented yet"
}

// TODO updateOne one update document data
func (db *DB) updateOne(query string) (result string) {
	return "not implemented yet"
}

// TODO updateMany update document data
func (db *DB) updateMany(query string) (result string) {
	return "not implemented yet"
}

// end
