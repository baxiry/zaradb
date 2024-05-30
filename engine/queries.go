package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// data that matched
type matched struct {
	id   string
	data string
}

// deletes Many items
func (db *DB) deleteMany(query string) string {

	mtch := gjson.Get(query, "match").String()

	coll := gjson.Get(query, "collection").String()

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}

	record := ""
	id := "1"

	listMatch := []string{}

	for rows.Next() {
		err := rows.Scan(&id, &record)
		if err != nil {
			return err.Error() // TODO standaring errors
		}

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}
		if ok {
			println("id: ", id)
			listMatch = append(listMatch, id)
		}
	}
	rows.Close()

	for _, id := range listMatch {
		// TODO use where in (1,2,3) for speedup query
		stmt := `DELETE FROM ` + coll + ` WHERE rowid = ` + id + `;`
		_, err = db.db.Exec(stmt)
		if err != nil {
			fmt.Println("delete erroo", err.Error())
			return `{"error": "` + err.Error() + `"}`
		}
	}
	return "ok"
}

// TODO updateMany update document data
func (db *DB) updateMany(query string) (result string) {

	mtch := gjson.Get(query, "match").String()

	newObj := gjson.Get(query, "data").String()

	coll := gjson.Get(query, "collection").String()

	// updates exist value

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}

	record := ""
	id := "1"

	listMatch := []matched{}

	for rows.Next() {
		err := rows.Scan(&id, &record)
		if err != nil {
			return err.Error() // TODO standaring errors
		}

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}
		if ok {
			println("id: ", id)
			listMatch = append(listMatch, matched{id: id, data: record})
		}
	}
	rows.Close()

	for _, rec := range listMatch {

		newData := gjson.Get(`[`+rec.data+`,`+newObj+`]`, `@join`).Raw
		// update test set record = '{"_id":38}' where rowid in (36, 37,38,39,40);

		stmt = `UPDATE ` + coll + ` SET record  = '` + newData + `' WHERE rowid = ` + rec.id + ";"
		_, err = db.db.Exec(stmt)
		if err != nil {
			fmt.Println(err.Error())
			// TODO
			return `{"error": "` + err.Error() + `"}`
		}
	}

	return fmt.Sprint(len(listMatch)) + " items updated"
}

// TODO updateOne one update document data
func (db *DB) updateOne(query string) (result string) {

	mtch := gjson.Get(query, "match").String()

	newObj := gjson.Get(query, "data").String()

	coll := gjson.Get(query, "collection").String()

	// updates exist value

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}

	record := ""
	id := "1"
	for rows.Next() {
		err := rows.Scan(&id, &record)
		if err != nil {
			return err.Error() // TODO standaring errors
		}

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}
		if ok {
			break
		}

	}
	rows.Close()

	newData := gjson.Get(`[`+record+`,`+newObj+`]`, `@join`).Raw
	// update test set record = '{"_id":12,"name":"joha","age":13}' where rowid = 39;
	stmt = `UPDATE ` + coll + ` SET record  = '` + newData + `' WHERE rowid = ` + id + ";"
	_, err = db.db.Exec(stmt)
	if err != nil {
		fmt.Println(err.Error())
		// TODO
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"update:": "done"}`
}

// Update update document data
func (db *DB) updateById(query string) (result string) {
	oldObj := db.findById(query)

	id := gjson.Get(query, "_id").String()
	newObj := gjson.Get(query, "data").String()
	coll := gjson.Get(query, "collection").String()

	newData := gjson.Get(`[`+oldObj+`,`+newObj+`]`, `@join`).Raw

	// updates exist value

	_ = `update test set record = '{"_id":12,"name":"joha","age":13}' where rowid = 39;`
	stmt := `UPDATE ` + coll + ` SET record  = '` + newData + `' WHERE rowid = ` + id + ";"
	_, err := db.db.Exec(stmt)
	if err != nil {
		return err.Error()
	}

	return "update done"
}

// Find finds any obs match creteria.
func (db *DB) findMany(query string) (res string) {

	// TODO parse hol qury one time
	coll := gjson.Get(query, "collection").String()
	if coll == "" {
		return `{"error":"forgot collection name "}`
	}

	mtch := gjson.Get(query, "match").String()

	if mtch == "" {
		mtch = "{}"
	}

	skip := gjson.Get(query, "skip").Int()
	limit := gjson.Get(query, "limit").Int()
	if limit == 0 {
		limit = 100 // what is default setting ?
	}

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	record := ""
	listData := make([]string, 0)
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
			return err.Error() // TODO standaring errors
		}

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}

		if ok {
			// aggrigate here
			listData = append(listData, record)
			limit--
		}
	}

	// TODO aggrigate here

	// remove|rename some fields
	flds := gjson.Get(query, "fields")
	listData = fields(listData, flds)

	records := "["

	for i := 0; i < len(listData); i++ {
		records += listData[i] + ","
	}

	if len(records) == 1 {
		return records + "]"
	}

	return records[:len(records)-1] + "]"
}

// Finds first obj match creteria.
func (db *DB) findOne(query string) (res string) {
	coll := gjson.Get(query, "collection").String()

	mtch := gjson.Get(query, "match").String()
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
		b, err := match(mtch, record)
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
	mtch := gjson.Get(query, "match").String()

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
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

		b, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}
		if b {
			break
		}
	}
	// should close here
	rows.Close()

	_, err = db.db.Exec(`delete from ` + coll + ` where rowid = ` + rowid) // fast
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	return `{"result":"` + rowid + ` removed"}`

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

	_, err := db.db.Exec(sql)
	if err != nil {
		return `{"error": "internal error"}` // + err.Error()
	}
	return `{"aknowlge": "row ` + id + ` deleted"}`
}

// end
