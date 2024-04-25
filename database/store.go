package database

import (
	"database/sql"
	"fmt"
	"strings"

	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db     *sql.DB
	lastid map[string]int64
}

func NewDB(dbName string) *DB {
	lastid := make(map[string]int64)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}

	// Query the sqlite_master table to get table names
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
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
		lid, _ := getLastId(db, tableName)
		lastid[tableName] = lid

		fmt.Println("Table Name:", tableName)
	}

	return &DB{db: db, lastid: lastid}
}

func getLastId(db *sql.DB, table string) (int64, error) {
	stmt := fmt.Sprintf(`SELECT rowid FROM %s ORDER BY ROWID DESC LIMIT 1`, table)
	res, err := db.Query(stmt)
	if err != nil {
		return 0, err
	}
	defer res.Close()

	var lastid int64

	if !res.Next() {
		return 0, nil
	}

	err = res.Scan(&lastid)
	if err != nil {
		return 0, err
	}

	return lastid, err
}

func (db *DB) createCollection(coll string) error {
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (obj json);`, coll)
	_, err := db.db.Exec(query)
	if err != nil {
		return err
	}
	li, err := getLastId(db.db, coll)
	if err != nil {
		return err
	}
	db.lastid[coll] = li
	return nil
}

func (db *DB) insert(coll, obj string) error {
	db.lastid[coll]++

	d := strings.TrimLeft(obj, " ")
	data := `{"_id":` + fmt.Sprint(db.lastid[coll]) + ", " + d[1:]

	fmt.Println("set json: ", data)

	println(db.lastid[coll], data)

	_, err := db.db.Exec(`insert into ` + coll + `(obj) values('` + data + `');`) // fast
	if err != nil {
		db.lastid[coll]--
	}
	return err
}

// =====================================================

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
