package kvlite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Path string
}

type index struct {
	// location format is :
	// "i <id> <at> <size> <page-name>"
	// "i 0 199 45 0"

	coll string
	page int
	id   int64
	at   int64
	size int
}

type Database struct {
	name   string
	page   int
	index  int64
	lastat int64
	// int for file name will be emprove speed a lettel bit
	pages  map[string]*os.File
	indexs []index
	afile  string // active file
	path   string
}

// Set inserts new or update exist value
func (db *Database) Set(key int, coll, value string) {

	size := len(value)

	// TODO use string builder to reduce memory consomption
	location := "\ni " + fmt.Sprint(key) + " " + fmt.Sprint(db.lastat) + " " + fmt.Sprint(size) + " 0" + coll + "\n"

	db.pages[db.afile].Write([]byte(value + location))

	// TODO assigne by index myght emprove performece a little bit.
	// db.indexs[key] = index{at: db.lastat, size: size, page: db.page}
	db.indexs = append(db.indexs, index{
		at:   db.lastat,
		size: size,
		coll: coll,
		page: db.page})

	db.lastat += int64(size + len(location))
}

// Get data by key
func (db *Database) Get(key int) string {

	// location format is :
	// "i <key> <at> <size> <page-name>"
	// "i 0 199 45 0"
	index := db.indexs[key]

	buffer := make([]byte, index.size)

	db.pages[db.path+fmt.Sprint(index.page)].ReadAt(buffer, index.at)

	// TODO make  reuse value will be improve mem & reduce gc
	// db.value ?!

	return string(buffer)
}

// Open initialaze db pages
func Open(path string) *Database {

	db := &Database{}

	db.indexs = db.reIndex()
	db.pages = make(map[string]*os.File)
	afile := "0" // active file
	db.path = path

	if db.path == "" {
		//path, _ = os.Getwd()
		db.path = "mok/"

		err := os.Mkdir(db.path, 0744)
		check("Mkdir ", err)

		db.afile = db.path + afile // active file

		file, err := os.OpenFile(db.afile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		check("when open file", err)

		fmt.Println("file active is : ", file.Name())
		db.pages[db.afile] = file

		// complet db initalaze

		return db
	}

	err := os.Mkdir(db.path, 0744)
	check("Mkdir ", err)

	db.afile = db.path + afile // active file

	file, err := os.OpenFile(db.afile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	check("when open file", err)
	file.Close()

	files, err := os.ReadDir(db.path)
	check("ReadDir ", err)

	for _, f := range files {

		dpage := db.path + f.Name()

		file, err := os.OpenFile(dpage, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		check("", err)

		// TODO int as file name
		db.pages[dpage] = file

		db.afile = dpage
	}

	return db
}

// rebuilds indexs
func (db *Database) reIndex() (indexs []index) {
	// Read the entire file into a byte slice

	indexs = make([]index, 0)

	for f := range db.pages {
		fileContent, err := os.ReadFile(f)
		check("", err)

		// Split the byte slice into lines using the newline character as a delimiter
		lines := strings.Split(string(fileContent), "\n")

		// Process each line
		for _, line := range lines {
			if len(line) == 0 {
				return
			}
			if line[0] == 'i' {

				pos := strings.Fields(line)
				at, _ := strconv.Atoi(pos[2])
				size, _ := strconv.Atoi(pos[3])
				id, _ := strconv.Atoi(pos[1])

				// TODO load pages and collections
				indexs[id] = index{at: int64(at), size: size}
			}
		}
	}
	return indexs
}

// Close db
func (db *Database) Close() {
	for _, f := range db.pages {
		f.Close()
	}
}

// error
func check(hint string, err error) {
	if err != nil {
		fmt.Println(hint, err)
		//return
	}
}

func (db *Database) ShowIndexs() {

	for k, v := range db.indexs {
		fmt.Println(k, v)
	}

	fmt.Println("len indexs : ", len(db.indexs))
}
