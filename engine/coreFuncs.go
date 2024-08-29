package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// data that matched
type matched struct {
	id   string
	data string
}

// ...
func getData(query gjson.Result) (data []string, err error) {
	coll := query.Get("collection").Str
	if coll == "" {
		return nil, fmt.Errorf(`{"error":"forgot collection name "}`)
	}

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		limit = 1000 // what is default setting ?
	}

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	record := ""

	isMatch := query.Get("match")

	for rows.Next() {

		if limit == 0 {
			break
		}

		record = ""
		err := rows.Scan(&record)
		if err != nil {
			return nil, err
		}

		ok, err := match(isMatch, record)
		if err != nil {
			return nil, err
		}

		if ok {
			if skip != 0 {
				skip--
				continue
			}
			data = append(data, record)
			limit--
		}

	}
	return data, nil
}

// deletes Many items
func (db *DB) deleteMany(query gjson.Result) string {

	mtch := query.Get("match")

	coll := query.Get("collection").Str

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

	llist := len(listMatch)
	for i := 0; i < llist; i++ {
		// TODO use where in (1,2,3) for speedup query
		stmt := `DELETE FROM ` + coll + ` WHERE rowid = ` + listMatch[i] + `;`
		_, err = db.db.Exec(stmt)
		if err != nil {
			fmt.Println("delete erro", err.Error())
			return `{"error": "` + err.Error() + `"}`
		}
	}
	return strconv.Itoa(llist) + "items has removed"
}

// TODO updateMany update document data
func (db *DB) updateMany(query gjson.Result) (result string) {

	mtch := query.Get("match")

	newObj := query.Get("data").Str

	coll := query.Get("collection").Str

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

	return strconv.Itoa(len(listMatch)) + " items updated"
}

// TODO updateOne updates one  document data
func (db *DB) updateOne(query gjson.Result) (result string) {

	mtch := query.Get("match")

	newObj := query.Get("data").Str

	coll := query.Get("collection").Str

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
func (db *DB) updateById(query gjson.Result) (result string) {

	oldObj := db.findById(query)
	fmt.Println("query ", query)

	id := query.Get("_id").String()
	if id == "" {
		return `{"error": "forget _id"}`
	}
	newObj := query.Get("data").Raw

	coll := query.Get("collection").Str

	newData := gjson.Get(`[`+oldObj+`,`+newObj+`]`, `@join`).Raw
	fmt.Println()
	fmt.Println("old obj ", oldObj)
	fmt.Println("new obj ", newObj)
	fmt.Println("new data ", newData)

	// example `update test set record = '{"_id":12,"name":"joha","age":13}' where rowid = 12;`

	stmt := `UPDATE ` + coll + ` SET record  = '` + newData + `' WHERE rowid = ` + id + ";"
	_, err := db.db.Exec(stmt)
	if err != nil {
		return err.Error()
	}

	return id + " updated"
}

// Find finds any obs match creteria.
func (db *DB) findMany(query gjson.Result) (res string) {

	listData, err := getData(query)
	if err != nil {
		return err.Error()
	}

	// order :
	order := query.Get("sort") //.Str

	if order.Exists() {
		listData = orderBy(order, 0, listData)
	}

	// TODO aggrigate here

	// remove or rename some fields
	flds := query.Get("fields")
	listData = reFields(listData, flds)

	records := "["

	for i := 0; i < len(listData); i++ {
		records += listData[i] + ","
	}

	ln := len(records)
	if ln == 1 {
		return "[]"
	}

	return records[:ln-1] + "]"
}

// Finds first obj match creteria.
func (db *DB) findOne(query gjson.Result) (res string) {
	coll := query.Get("collection").Str
	skip := query.Get("skip").Int()

	// TODO are skyp useful here ?

	stmt := `select record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}
	defer rows.Close()

	mtch := query.Get("match")
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

		ok, err := match(mtch, record)
		if err != nil {
			return err.Error()
		}

		if ok {

			return record
		}
	}

	return `{"result":0}`
}

// delete
func (db *DB) deleteOne(query gjson.Result) string {

	coll := query.Get("collection").Str

	stmt := `select rowid, record from ` + coll

	rows, err := db.db.Query(stmt)
	if err != nil {
		return err.Error()
	}

	rowid := "0"
	record := ""
	mtch := query.Get("match")
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

	_, err = db.db.Exec(`delete from ` + coll + ` where rowid = ` + rowid) // + is fast
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	return `{"result":"_id:` + rowid + ` removed"}`

}

// Finds first obj match creteria.
func (db *DB) findById(query gjson.Result) (res string) {
	coll := query.Get("collection").Str

	id := query.Get("_id").String()

	stmt := `select record from ` + coll + ` where rowid = ` + id

	rows, err := db.db.Query(stmt)
	if err != nil {
		if id == "" {
			return `{Error:"_id is required"}`
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

// Insert One
func (db *DB) insertOne(query gjson.Result) (res string) {
	coll := query.Get("collection").Str
	data := query.Get("data").String() // .Str not works with json obj
	fmt.Println(coll, "\n", data)

	err := db.insert(coll, data)
	if err != nil {
		//db.lastid[coll] = lid
		if strings.Contains(err.Error(), "no such table") {
			err = db.CreateCollection(coll)
			if err != nil {
				return err.Error()
			}
			err = db.insert(coll, data)
			return "inser done"
		}
		return err.Error()
	}

	return "inser done"
}

// InsertMany inserts list of object at one time
func (db *DB) insertMany(query gjson.Result) (res string) {
	coll := query.Get("collection").Str
	data := query.Get("data").Array()

	lid := db.lastid[coll]
	strData := ""
	for _, obj := range data {
		db.lastid[coll]++
		// strconv for per
		strData += `('{"_id":` + fmt.Sprint(db.lastid[coll]) + ", " + obj.String()[1:] + `'),`
	}

	_, err := db.db.Exec(`insert into ` + coll + `(record) values` + strData[:len(strData)-1]) // `+` is fast
	if err != nil {
		db.lastid[coll] = lid
		if strings.Contains(err.Error(), "no such table") {
			err = db.CreateCollection(coll)
			if err != nil {
				return err.Error()
			}
			return db.insertMany(query)
		}
		return err.Error()
	}

	return "inserted"
}

// delete by id
func (db *DB) deleteById(query gjson.Result) string {
	id := query.Get("_id").String()
	if id == "" {
		return `{"error": "_id is required"}`
	}
	coll := query.Get("collection").Str
	if coll == "" {
		return `{"error": "collection is required"}`
	}

	sql := `delete from ` + coll + ` where rowid = ` + id

	_, err := db.db.Exec(sql)
	if err != nil {
		return `{"error": "internal error"}` // + err.Error()
	}
	return `{"aknowlge": "row ` + id + ` deleted"}`
}

// end
