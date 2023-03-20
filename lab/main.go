package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func mainc() {

	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`DROP TABLE IF EXISTS cars; CREATE TABLE cars(id INTEGER PRIMARY KEY, value TEXT);`)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(`INSERT INTO cars(value) VALUES(?);`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Exec("hello how are you if yor there ??ok " + strconv.Itoa(2))

	fmt.Println("table cars created")
}
