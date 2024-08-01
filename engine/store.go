package engine

import (
	"database/sql"
	"fmt"
	"strings"

	//_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
	_ "modernc.org/sqlite"
)

type DB struct {
	db     *sql.DB
	lastid map[string]int64
}

// var lastid = make(map[string]int64, 0)
var db *DB

func NewDB(dbName string) *DB {

	newdb, err := sql.Open("sqlite", dbName) // sqlite3 whith mattn lib
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

		fmt.Println("Table Name:", tableName)
	}
	return db
}

// InsertMany inserts list of object at one time
func (db *DB) insertMany(query gjson.Result) (res string) {
	coll := query.Get("collection").Str
	data := query.Get("data").Array()

	//d := strings.TrimLeft(obj, " ")
	// if len(d) == 0 {return fmt.Sprintf("len data is 0 %s\n", d)}

	lid := db.lastid[coll]
	strData := ""
	for _, obj := range data {
		db.lastid[coll]++
		// strconv for per
		strData += `('{"_id":` + fmt.Sprint(db.lastid[coll]) + ", " + obj.String()[1:] + `'),`
	}

	//fmt.Println("bulk data:  ", strData[:len(strData)-1])

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
