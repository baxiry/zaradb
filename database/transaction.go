package database

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func transaction(query string) string {

	actions := gjson.Get(query, "transaction").Array()
	start := "t " + str(len(actions)) + "\n"
	_ = start
	for k, v := range actions {

		fmt.Println(k, v)
	}
	return "actions done"
}

/*
// updates exist value
func (db *Database) tupdate(id int, coll, value string) string {
	if id > db.lastId {
		return "Id not exists"
	}

	if db.indexs[id].coll != coll && db.indexs[id].coll != "" {
		return "coll not match"
	}
	if db.indexs[id].at == 0 {
		//return "item not exists"
	}

	size := len(value)
	page := " 0 "

	// TODO use string builder to reduce memory consomption
	location := "\ni " + str(id) + " " + str(db.lat) + " " + str(size) + page + coll + "\n"

	db.pages[db.activeFile].Write([]byte(value + location))

	db.indexs[id] = index{at: db.lat, size: size, coll: coll, page: db.page}

	db.lat += int64(size + len(location))

	return "done"
}
*/
