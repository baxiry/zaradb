package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	db, err := sql.Open("sqlite3", "test.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rows, err := db.Query("SELECT * FROM cars")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var name string
		var price int

		err = rows.Scan(&id, &name, &price)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d %s %d\n", id, name, price)
	}
}
