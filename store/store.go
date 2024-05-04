package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
)

type DB struct {
	db *sql.DB
}

var db *DB

var lastid = make(map[string]int64, 0)

/*
INSERT INTO your_table_name (column)
VALUES (value1), (value2), (value3);
*/
// Insert
func (db *DB) insertMany(query string) (res string) {
	coll := gjson.Get(query, "collection").String()
	data := gjson.Get(query, "data").Array()

	//d := strings.TrimLeft(obj, " ")
	// if len(d) == 0 {return fmt.Sprintf("len data is 0 %s\n", d)}

	lid := lastid[coll]
	strData := ""
	for _, obj := range data {
		lastid[coll]++
		strData += `('{"_id":` + fmt.Sprint(lastid[coll]) + ", " + obj.String()[1:] + `'),`
	}

	fmt.Println("bulk data:  ", strData[:len(strData)-1])

	_, err := db.db.Exec(`insert into ` + coll + `(record) values` + strData[:len(strData)-1]) // fast
	if err != nil {
		lastid[coll] = lid
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

	lastid[collection]++
	data := `{"_id":` + fmt.Sprint(lastid[collection]) + ", " + d[1:]
	_, err := db.db.Exec(`insert into ` + collection + `(record) values('` + data + `');`) // fast
	if err != nil {
		println(err)
		lastid[collection]--
		return err
	}

	return nil
}

func NewDB(dbName string) *DB {

	lastid = make(map[string]int64, 0)

	newdb, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

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

		lastid[tableName] = lid

		fmt.Println("Table Name:", tableName)
	}
	db = &DB{db: newdb}

	return db
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
	lastid[collection] = lid
	return nil
}

// =====================================================

// Close db
func (db *DB) Close() {
	db.db.Close()
}

// deletes exist value
func (db *DB) Delete(id int, coll string) string {
	return "not implement yet!"
}

// updates exist value
func (db *DB) Update(id int, coll, value string) string {
	return "not implemented yet"
}

// Get data by id
func (db *DB) Get(id int, coll string) string {
	return "get not emplement yet"
}

func Run(path string) {
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
