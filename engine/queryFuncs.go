package engine

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// TODO  use strings.Builder for string concatunation

// Insert One
func (store *Store) insertOne(query gjson.Result) (res string) {

	coll := query.Get("collection").Str
	if coll == "" {
		return `{"error":"forgot collection"}`
	}

	data := query.Get("data").String() // .Str not works with json obj
	if data == "" {
		return `{"error":"forgot data"}`
	}

	store.lastids[coll]++
	key := strconv.Itoa(int(store.lastids[coll]))
	data = `{"_id":` + key + ", " + data[1:] // strings.Builder

	err := store.Put(coll, data)
	if err != nil {
		if strings.Contains(err.Error(), "no such table:") {
			store.createCollection(query)
			err = store.Put(coll, data)
		} else {
			return `{"ak":"error:` + err.Error() + `"}`
		}
	}
	return `{"ak":"insert new obj with:` + key + ` success"}`
}

// InsertMany inserts list of object at one time
func (store *Store) insertMany(query gjson.Result) (res string) {
	coll := query.Get(collection).Str
	data := query.Get("data").Array()

	for _, obj := range data {
		store.lastids[coll]++
		key := strconv.Itoa(int(store.lastids[coll]))
		obj := `{"_id":` + key + `,` + obj.String()[1:]

		err := store.Put(coll, obj)

		if err != nil {
			if strings.Contains(err.Error(), "no such table:") {
				s.createCollection(query)
				err = s.Put(coll, obj)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				return `{"error":"` + err.Error() + `"}`
			}
		}
	}

	return `{"ak":"insertMany Done"}`
}

// Finds first obj match creteria.
func (s *Store) findOne(query gjson.Result) (rowid string, res string) {
	coll := query.Get(collection).Str
	skip := query.Get("skip").Int()
	isMatch := query.Get("match")

	rows, err := s.db.Query("SELECT rowid, obj FROM " + coll)
	if err != nil {
		return rowid, `{"error": "` + err.Error() + `"}`
	}

	var obj string
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&rowid, &obj)
		if err != nil {
			fmt.Println("err when Scan: ", err)
			continue
		}
		fmt.Println("rowid & obj is : ", rowid, obj)

		ok, err := match(isMatch, obj)
		if err != nil {
			return rowid, `{"error":"` + err.Error() + `"}`
		}

		if ok {
			if skip != 0 {
				skip--
				continue
			}
			res = obj
			break
		}
	}

	if res != "" {
		return rowid, res
	}

	return rowid, `{"status":"nothing match"}`
}

// Find finds any object match creteria.
func (s *Store) findMany(query gjson.Result) (res string) {

	listData, err := s.Get(query)
	if err != nil {
		return `{"error":"` + err.Error() + `"}`
	}

	// order :
	srt := query.Get("sort")

	if srt.Exists() {
		listData = order(listData, srt)
	}

	// remove or rename some fields
	flds := query.Get("fields")
	listData = reFields(listData, flds)

	records := "["
	for i := 0; i < len(listData); i++ {
		// Todo use strings.Builder
		records += listData[i] + ","
	}

	ln := len(records)
	if ln == 1 {
		return "[]"
	}

	// Todo use strings.Builder
	return records[:ln-1] + "]"
}

// TODO updateOne updates one  document data
func (s *Store) updateOne(query gjson.Result) (result string) {

	rowid, oldObj := s.findOne(query)
	if rowid == "" {
		return `{"error":"nothing to update"}`
	}

	newObj := query.Get("data").Raw
	newData := gjson.Get(`[`+oldObj+`,`+newObj+`]`, `@join`).Raw

	coll := query.Get(collection).Str

	stmt, err := s.db.Prepare("UPDATE " + coll + " SET obj = ? where rowid= ?;") // strings.Builder
	if err != nil {
		return err.Error()
	}

	_, err = stmt.Exec(newData, rowid)
	if err != nil {
		return `{"ak", "` + err.Error() + `"}` // TODO err
	}

	return `{"ak": "update: done"}`
}

// Finds first obj match creteria.
func (s *Store) findById(query gjson.Result) (res string) {
	coll := query.Get(collection).Str

	id := query.Get("_id").String()
	if id == "0" {
		return "_id forgoten"
	}

	row := s.db.QueryRow("SELECT obj FROM " + coll + " where rowid =" + id)
	err := row.Scan(&res)
	if err != nil {
		if strings.Contains(err.Error(), "no such table:") {
			return `{"error": "collection does not exist"}`
		}
		return `{"error": "_id does not exist"}`
	}

	return res
}

// TODO updateMany update document data
func (s *Store) updateMany(query gjson.Result) (result string) {

	fmt.Println("update Many")
	isMatch := query.Get("match")

	newObj := query.Get("data").Raw

	coll := query.Get("collection").Str

	rows, err := s.db.Query("select * from " + coll)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	for rows.Next() {
		var value string
		ok, err := match(isMatch, string(value))
		if err != nil {
			fmt.Println("updateMany err", err)
			continue
		}

		if ok {
			newData := gjson.Get(`[`+string(value)+`,`+newObj+`]`, `@join`).Raw
			_, err = s.db.Exec("Update "+coll+" set obj = "+newData, "where id  = 1")
			if err != nil {
				return `{"error": "` + err.Error() + `"}`
			}
		}

		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}

	}

	return "many items updated succesfully"
}

// delete one item
func (s *Store) deleteOne(query gjson.Result) string {

	coll := query.Get("collection").Str

	matchPattren := query.Get("match")

	rows, err := s.db.Query("select rowid from " + coll)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	var rowid string

	for rows.Next() {
		rows.Scan(&rowid)
		ok, err := match(matchPattren, rowid)
		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}

		if ok {
			_, err = s.db.Exec("delete from " + coll + "where rowid=" + rowid)
			if err != nil {
				return `{"error": "` + err.Error() + `"}`
			}

		}
	}

	return `{"result":"_id:` + rowid + ` deleted"}`
}

// deletes Many items
func (s *Store) deleteMany(query gjson.Result) string {

	isMatch := query.Get("match")

	coll := query.Get("collection").Str

	rows, err := s.db.Query("select * from " + coll)
	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	for rows.Next() {
		var obj string
		var rowid string
		rows.Scan(&rowid, &obj)

		ok, err := match(isMatch, obj)
		if err != nil {
			fmt.Println("updateMany err", err)
			continue
		}

		if ok {
			_, err = s.db.Exec("DELETE FROM "+coll+" WHERE id = ?", rowid)
			if err != nil {
				return `{"error": "` + err.Error() + `"}`
			}
		}

		if err != nil {
			return `{"error": "` + err.Error() + `"}`
		}
	}

	return " many items has removed successfully"
}

// Update update document data
func (s *Store) updateById(query gjson.Result) (result string) {

	id := query.Get("_id").Str
	if id == "" {
		return `{"error": "forget _id"}`
	}
	newObj := query.Get("data").Raw
	if id == "" {
		return `{"error": "forget data"}`
	}

	coll := query.Get("collection").Str
	var rowid string
	var oldObj string

	row := s.db.QueryRow("select * from " + coll)

	row.Scan(&rowid, &oldObj)

	newData := gjson.Get(`[`+oldObj+`,`+newObj+`]`, `@join`).Raw

	_, err := s.db.Exec("Update "+coll+"set obj ="+newData+"where rowid=", rowid)

	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return id + " updated"
}

// delete by id
func (s *Store) deleteById(query gjson.Result) string {

	id := query.Get("_id").Str
	if id == "" {
		return `{"error": "forget _id"}`
	}

	coll := query.Get("collection").Str

	_, err := s.db.Exec("delete from " + coll + "where rowid = " + id)

	if err != nil {
		return `{"error": "` + err.Error() + `"}`
	}

	return `{"aknowlge": "row ` + id + ` deleted"}`
}

// end
