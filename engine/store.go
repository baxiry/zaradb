package engine

import (
	"database/sql"
	"errors"
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
func (db *DB) insertMany(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").Array()

	//d := strings.TrimLeft(obj, " ")
	// if len(d) == 0 {return fmt.Sprintf("len data is 0 %s\n", d)}

	lid := db.lastid[coll]
	strData := ""
	for _, obj := range data {
		db.lastid[coll]++
		strData += `('{"_id":` + fmt.Sprint(db.lastid[coll]) + ", " + obj.String()[1:] + `'),`
	}

	fmt.Println("bulk data:  ", strData[:len(strData)-1])

	_, err := db.db.Exec(`insert into ` + coll + `(record) values` + strData[:len(strData)-1]) // fast
	if err != nil {
		db.lastid[coll] = lid
		return err.Error()
	}

	return "in progress"
}

// insert new record
func (db *DB) insert(collection, obj string) error {
	d := strings.TrimLeft(obj, " ")
	if len(d) < 2 {
		return fmt.Errorf("len data is 0 %s\n", d)
	}

	db.lastid[collection]++
	data := `{"_id":` + fmt.Sprint(db.lastid[collection]) + ", " + d[1:]
	_, err := db.db.Exec(`insert into ` + collection + `(record) values('` + data + `');`) // fast
	if err != nil {
		//println(err)
		db.lastid[collection]--
		return err
	}

	return nil
}

func (db *DB) CreateCollection(collection string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (record json);`, collection)
	_, err := db.db.Exec(query)
	if err != nil {
		return err
	}
	lid, err := getLastId(db.db, collection)
	if err != nil {
		return err
	}
	db.lastid[collection] = lid
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

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
