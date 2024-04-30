package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

var lastid = make(map[string]int64, 0)

// insert new record
func (db *DB) insert(collection, obj string) error {
	lastid[collection]++

	d := strings.TrimLeft(obj, " ")
	if len(d) == 0 {
		return fmt.Errorf("len data is 0 %s\n", d)
	}

	data := `{"_id":` + fmt.Sprint(lastid[collection]) + ", " + d[1:]

	fmt.Println("set json: ", data)

	println(lastid[collection], data)

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

	return &DB{db: newdb}
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
