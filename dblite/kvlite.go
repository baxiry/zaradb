package dblite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// max items per page
const MaxItems = 10

type Config struct {
	Path string
}

var db *Database

func Run(path string) *Database {
	db = Open(path)
	return db
}

type Database struct {
	page    int
	Lid     int
	lindexs int
	lat     int64 // last at

	// TODO page []*os.File
	pages map[string]*os.File

	//indexs []index
	indexs map[int]index
	afile  string // active file
	path   string
}

type index struct {
	// location format is :
	// "i <id> <at> <size> <page-name> <coll>"
	// "i 0 199 45 0 users"
	//id   int

	at   int64
	size int
	page int
	coll string
}

// deletes exist value
func (db *Database) Delete(id int, coll string) string {

	if id > db.Lid {
		return "Id not exists"
	}

	// indx, ok := db.indexs[id]; if !ok { return "no data to delete"	}

	indx := db.indexs[id]

	if indx.size == 0 {
		return "no data to delete"
	}

	if indx.coll != coll {
		return "coll wrong"
	}

	location := "d " + str(id) + "\n"

	db.pages[db.afile].Write([]byte(location))

	//delete(db.indexs, id)
	db.indexs[id] = index{}
	db.lat += int64(len(location))

	return "delete success!"
}

var str = fmt.Sprint

// updates exist value
func (db *Database) Update(id int, coll, value string) string {
	if id > db.Lid {
		return "Id not exists"
	}

	if db.indexs[id].coll != coll && db.indexs[id].coll != "" {
		return "coll not match"
	}
	if db.indexs[id].at == 0 {
		//return "item not exists"
	}

	size := len(value)
	page := " 0 "

	// TODO use string builder to reduce memory consomption
	location := "\ni " + str(id) + " " + str(db.lat) + " " + str(size) + page + coll + "\n"

	db.pages[db.afile].Write([]byte(value + location))

	db.indexs[id] = index{at: db.lat, size: size, coll: coll, page: db.page}

	db.lat += int64(size + len(location))

	return "done"
}

// last primary index
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

	db.Lid++

	size := len(value)
	page := " 0 "

	// TODO use string builder to reduce memory consomption
	location := "\ni " + str(db.Lid) + " " + str(db.lat) + " " + str(size) + page + coll + "\n"

	db.pages[db.afile].Write([]byte(value + location))

	db.indexs[db.Lid] = index{at: db.lat, size: size, coll: coll, page: db.page}

	db.lat += int64(size + len(location))
}

// Get data by key
func (db *Database) Get(id int, coll string) string {

	// location format is :
	// "i <id> <at> <size> <page> <coll>"

	// "i 1 0 33 0 users"
	if id > db.Lid {
		return "Id not exists"
	}

	index := db.indexs[id]
	if index.size == 0 {
		return "not exist"
	}

	//if index.coll != coll {return "coll not match"}

	buffer := make([]byte, index.size)

	// TODO cange page's type to list to improve cpu
	db.pages[db.path+fmt.Sprint(index.page)].ReadAt(buffer, index.at)

	// TODO make  reuse value will be improve mem & reduce gc
	// db.value ?!

	return string(buffer)
}

// rebuilds indexs
// func (db *Database) reIndex() (indexs map[int]index) {
func (db *Database) reIndex() (indexs map[int]index) {

	indexs = make(map[int]index, MaxItems)

	pages, err := os.ReadDir(db.path)
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
				continue
			}
			if line[0] == 'i' {
				pos := strings.Fields(line)

				// "i 1 0 33 0 users"
				id, _ := strconv.Atoi(pos[1])
				if id > db.Lid {
					db.Lid = id
				}
				at, _ := strconv.Atoi(pos[2])
				size, _ := strconv.Atoi(pos[3])
				indexs[id] = index{at: int64(at), size: size, coll: pos[5]}

			} else if line[0] == 'd' {
				pos := strings.Fields(line)
				id, _ := strconv.Atoi(pos[1])

				indexs[id] = index{}
				//delete(indexs, id)
			}
		}
	}

	db.lastAt()

	fmt.Println("last id : ", db.Lid)
	return indexs
}

// Open initialaze db pages
func Open(path string) *Database {

	db = &Database{}

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
		db.lindexs = len(db.indexs)

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

	// TODO we need reIndex wen server crushed. not in normal stopt
	// for devloping mod i use reIndex fot testing
	db.indexs = db.reIndex()

	db.lindexs = len(db.indexs)
	return db
}

func (db *Database) saveIndexs() {

	/*
		// TODO save index for fast start if sever stoptd greacefully
		// if not then we new rebuild indexes from data

		for k, v := range db.indexs {
			//file.WriteString(fmt.Sprintf("%d %v\n", k, v))
			fmt.Printf("%d, %v\n", k, v)
		}

			file, _ := os.Create(db.path + "indexs")
			os.Create(db.path + "Done")
	*/

}

// Close db
func (db *Database) Close() {
	for _, f := range db.pages {
		f.Close()
	}

	// TODO
	db.saveIndexs()
}

// error
func check(hint string, err error) {
	if err != nil {
		fmt.Println(hint, err)
		//return
	}
}
