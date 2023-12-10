package kvlite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var str = fmt.Sprint

type Config struct {
	Path string
}

type index struct {
	// location format is :
	// "i <id> <at> <size> <page-name> <coll>"
	// "i 0 199 45 0 users"
	coll string
	page int
	// id   int64
	at   int64
	size int
}

type Database struct {
	page int
	Lid  int
	lat  int64 // last at

	// int for file name will be emprove speed a lettel bit
	pages map[string]*os.File
	//indexs map[int]index
	indexs map[int]index
	afile  string // active file
	path   string
}

// deletes exist value
func (db *Database) Delete(id int, coll string) {

	indx, ok := db.indexs[id]
	if !ok {
		println("no data to delete")
		return
	}
	if indx.coll != coll {
		fmt.Println("======== no match", indx.coll, ",", coll)
		return
	}

	location := "d " + str(id) + "\n"

	db.pages[db.afile].Write([]byte(location))

	delete(db.indexs, id)

	db.lat += int64(len(location))
}

// updates exist value
func (db *Database) Update(id int, coll, value string) {

	size := len(value)
	page := " 0 "

	// TODO use string builder to reduce memory consomption
	location := "\ni " + str(id) + " " + str(db.lat) + " " + str(size) + page + coll + "\n"

	db.pages[db.afile].Write([]byte(value + location))

	db.indexs[id] = index{at: db.lat, size: size, coll: coll, page: db.page}

	db.lat += int64(size + len(location))
}

func (db *Database) lastAt() {

	files, err := os.ReadDir(db.path)
	check("ReadDir ", err)

	for _, f := range files {
		dpage := db.path + f.Name()
		state, err := os.Stat(dpage)
		check("read state", err)
		db.lat += state.Size()
	}
}

// inserts new or update exist value
func (db *Database) Insert(coll, value string) {

	size := len(value)
	page := " 0 "

	// TODO use string builder to reduce memory consomption

	location := "\ni " + str(db.Lid) + " " + str(db.lat) + " " + str(size) + page + coll + "\n"

	db.pages[db.afile].Write([]byte(value + location))

	db.indexs[db.Lid] = index{at: db.lat, size: size, coll: coll, page: db.page}

	db.lat += int64(size + len(location))
	db.Lid++
}

// Get data by key
func (db *Database) Get(id int) string {

	// location format is :
	// "i <id> <at> <size> <page> <coll>"
	// "i 0 199 45 0 users"
	index, ok := db.indexs[id]
	if !ok {
		return "no val"
	}

	buffer := make([]byte, index.size)

	db.pages[db.path+fmt.Sprint(index.page)].ReadAt(buffer, index.at)

	// TODO make  reuse value will be improve mem & reduce gc
	// db.value ?!

	return string(buffer)
}

// rebuilds indexs
func (db *Database) reIndex() (indexs map[int]index) {

	// Read the entire file into a byte slice

	indexs = make(map[int]index)

	pages, err := os.ReadDir("db1")
	if err != nil {
		fmt.Println("read dir", err)
	}

	for _, f := range pages {
		fileContent, err := os.ReadFile(db.path + f.Name())
		check("range over files: ", err)

		// Split the byte slice into lines using the newline character as a delimiter
		lines := strings.Split(string(fileContent), "\n")

		// Process each line
		for _, line := range lines {
			if len(line) == 0 {
				fmt.Println("why this ? ")
				break
			}
			if line[0] == 'i' {

				pos := strings.Fields(line)
				at, _ := strconv.Atoi(pos[2])

				size, _ := strconv.Atoi(pos[3])
				id, _ := strconv.Atoi(pos[1])
				if id > db.Lid {
					db.Lid = id
				}

				// TODO elso add pages and collections
				indexs[id] = index{at: int64(at), size: size}
			} else if line[0] == 'd' {
				pos := strings.Fields(line)
				id, _ := strconv.Atoi(pos[1])
				delete(indexs, id)
			}
		}
	}

	if db.Lid != 0 {
		db.Lid++
	}
	//db.Lid = len(indexs)
	db.lastAt()

	fmt.Println("last id : ", db.Lid)
	return indexs
}

// Open initialaze db pages
func Open(path string) *Database {

	db := &Database{}

	//fmt.Println("last id is : ", db.lid)

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

		db.indexs = db.reIndex()

		// complet db initalaze
		return db
	}

	err := os.Mkdir(db.path, 0744)
	check("Mkdir", err)

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

	db.indexs = db.reIndex()
	return db
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
