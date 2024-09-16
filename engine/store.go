package engine

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type DB struct {
	db     *sql.DB
	lastid map[string]int64
}

// var lastid = make(map[string]int64, 0)
var db *DB

func NewDB(dbName string) *DB {

	newdb, err := sql.Open("sqlite", dbName)
	if err != nil {
		panic(err)
	}

	db = &DB{db: newdb}
	db.lastid = make(map[string]int64, 0)

	// Query the sqlite_master table to get table names
	rows, err := newdb.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterate through rows and print table names
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			panic(err)
		}
		lid, _ := getLastId(newdb, tableName)

		db.lastid[tableName] = lid

		//fmt.Println("Table Name:", tableName)
	}
	return db
}

// insert new record
func (db *DB) insert(collection, obj string) error {
	d := strings.TrimLeft(obj, " ")
	if len(d) < 2 {
		return fmt.Errorf("len data is 0 %s\n", d)
	}

	db.lastid[collection]++
	data := `{"_id":` + fmt.Sprint(db.lastid[collection]) + ", " + d[1:]
	fmt.Println("data: ", data)
	fmt.Println("coll: ", collection)
	// + s faster then format
	_, err := db.db.Exec(`insert into ` + collection + `(record) values('` + data + `');`)
	if err != nil {
		fmt.Println(err)
		db.lastid[collection]--
		return err
	}

	return nil
}

// Close db
func (db *DB) Close() {
	db.db.Close()
}

// error
func check(hint string, err error) {
	if err != nil {
		fmt.Println(hint, err)
		//return
	}
}

// end
