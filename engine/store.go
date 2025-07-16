package engine

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tidwall/gjson"
)

type Store struct {
	db      *sql.DB
	lastids map[string]uint64
}

var s = &Store{}

func NewDB(path string) *Store {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to ensure connectivity
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	s.db = db
	s.lastids = make(map[string]uint64, 1)

	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table';`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		tableName := ""
		rows.Scan(&tableName)
		rowid := db.QueryRow(`SELECT MAX(rowid) FROM ` + tableName)
		var lastid uint64
		rowid.Scan(&lastid)
		s.lastids[tableName] = lastid
		fmt.Printf("table: %s, lastid: %d\n", tableName, lastid)
	}

	return s
}

// insert
func (s *Store) Put(coll, val string) (err error) {
	_, err = s.db.Exec("INSERT INTO " + coll + `(obj) VALUES('` + val + `')`)
	return err
}

// getData fitch for data
func (s *Store) Get(query gjson.Result) (data []string, err error) {
	coll := query.Get("collection").Str
	if coll == "" {
		return nil, fmt.Errorf(`{"error":"forgot collection name "}`)
	}

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		// TODO what default number should be here ?
		limit = 100
	}

	isMatch := query.Get("match")

	rows, err := s.db.Query("select * from " + coll)
	if err != nil {
		fmt.Println("s.db.Query store.Get:", err)
		return nil, err
	}
	//defer rows.Close()

	var obj string

	for rows.Next() {
		if limit == 0 {
			break
		}

		rows.Scan(&obj)

		ok, err := match(isMatch, obj)
		if err != nil {
			return nil, err
		}

		if ok {
			if skip != 0 {
				skip--
				continue
			}
			data = append(data, obj)
			limit--
		}
	}
	rows.Close()

	return data, err
}

func (s *Store) Close() {
	s.db.Close()
}
