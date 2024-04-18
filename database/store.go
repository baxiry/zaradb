package database

import (
	"fmt"

	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) createCollection(collection string) error {

	query := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        doc JSON 
    );`, collection)

	_, err := db.db.Exec(query)
	if err != nil {
		fmt.Println("ERROR when creating table", err)
		return err
	}

	stmt := fmt.Sprintf(`insert into %s (doc) values('{"hi":"hello world"}');`, "test")
	_, err = db.db.Exec(stmt)
	if err != nil {
		fmt.Println("ERROR when inserting data into table", err)
	}
	_, err = db.db.Exec(query)
	return err
}

// =====================================================

// sprint alias
var str = fmt.Sprint

// Close db
func (db *DB) Close() {}

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

// Open initialaze db pages
func Open(path string) *DB {
	db := &DB{}
	return db
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
