package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

// read
func readSpicificLine(fn string, n int) (string, error) {
	if n < 1 {
		return "", fmt.Errorf("invalid request: line %d", n)
	}
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	bf := bufio.NewReader(f)
	var line string
	for lnum := 0; lnum < n; lnum++ {
		line, err = bf.ReadString('\n')
		if err == io.EOF {
			switch lnum {
			case 0:
				return "", errors.New("no lines in file")
			case 1:
				return "", errors.New("only 1 line")
			default:
				return "", fmt.Errorf("only %d lines", lnum)
			}
		}
		if err != nil {
			return "", err
		}
	}
	if line == "" {
		return "", fmt.Errorf("line %d empty", n)
	}
	return line, nil
}

// generate data to file
func generateLinse(n int) {
	s := "This is a string\n"

	for i := 0; i <= n; i++ {
		snum := strconv.Itoa(i)
		ioutil.WriteFile(snum+"_file.txt", []byte(snum+" "+s), 0644)
	}

}

// find a line
func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}
