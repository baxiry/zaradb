package dblite

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

// Root database folder
var RootPath string = rootPath() + slash
var pi = "pi" // primary index

type Database struct {
	Name        string
	Collections string
	Pages       map[string]*os.File
}

// NewCollection constracts List of files collection
func NewDatabase(name string) *Database {
	database := &Database{
		Name:        RootPath + name + slash,
		Collections: "test" + slash,
		Pages:       make(map[string]*os.File, 2),
	}
	return database
}

// creates new page and add it to Collections
func (db *Database) NewPage(id int) {
	// TODO
	indexFilePath := db.Name + db.Collections + "pi"

	filename, _, _ := GetIndex(db.Pages[indexFilePath], id)
	//	iLog.Println("GetIndex from :", indexFilePath)

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	path := filepath.Join(db.Name, db.Collections+strconv.Itoa(id))

	db.Pages[path] = file
	//iLog.Printf("new page is created with %s path\n", path)
}

// opnens all collection in Root database folder
func (db *Database) Open() {
	path := db.Name + db.Collections
	//	iLog.Println("opening database ", path)

	var err error
	var files []fs.DirEntry

	files, err = os.ReadDir(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(db.Name+db.Collections, 0744)
		if err != nil {
			eLog.Println("while mkDir", err)
		}
	}

	_, err = os.Stat(db.Name + db.Collections + pi)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(db.Name+db.Collections+pi, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			eLog.Println("when creating pi ", err)
			return
		}
		f.Close()
	}

	files, err = os.ReadDir(path)
	if err != nil {
		eLog.Printf("while reading dir %s, %v\n\n", path, err)
		return
	}

	//iLog.Printf("reading  %s\n", db.Name+db.Collections)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		page, err := os.OpenFile(path+file.Name(), os.O_RDWR, 0644) //
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+file.Name(), err)
			//break
		}
		// filepath.Join(path, file.Name())
		db.Pages[db.Name+db.Collections+file.Name()] = page

		//	iLog.Printf("%s is ready\n", file.Name())
	}
	if len(db.Pages) < 2 {
		page, err := os.OpenFile(path+"0", os.O_CREATE|os.O_RDWR, 0644) //
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+"0", err)
		}
		// filepath.Join(path, file.Name())
		db.Pages[db.Name+db.Collections+"0"] = page

	}
	// iLog.Println("length of db.Pages is : ", len(db.Pages))
}

// closes All collection
func (db *Database) Close() {
	for _, page := range db.Pages {
		page.Close()
		iLog.Printf("%s closed\n", page.Name())
	}
}
