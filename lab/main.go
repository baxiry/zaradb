package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	start := time.Now()

	args := os.Args
	s := string(args[1])
	input, _ := strconv.Atoi(s)

	//write(input)
	read(input)
	println("duration: ", time.Since(start).Milliseconds())
}

func write(input int) {
	db, err := sql.Open("sqlite3", "database.sql")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 1; i <= input; i++ {
		sts := fmt.Sprintf(`INSERT INTO cars(name) VALUES('Audi   %s');`, strconv.Itoa(i))

		_, err = db.Exec(sts)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func read(input int) {
	db, err := sql.Open("sqlite3", "database.sql")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stm, err := db.Prepare("SELECT name FROM cars WHERE id = ?")

	if err != nil {
		log.Fatal(err)
	}
	defer stm.Close()

	var name string
	var d string

	for i := 0; i < input; i++ {

		err = stm.QueryRow(3).Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		d += string(name[0])

	}
	fmt.Println("size data;", len(d))

	// inserData := `INSERT INTO cars(name) VALUES(' 52642');`
	// _, err = db.Exec(sts)

}
