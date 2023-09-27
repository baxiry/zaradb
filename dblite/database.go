package dblite

import (
	"os"
)

var pi = "pi" // primary index file

type Collection struct {
	at           int64
	primaryIndex int64
	// [[0,3],[3,8]]
	cachedIndexs [][2]int64
}

type Collections map[string]Collection

type Database struct {
	//Indexs       map[string]Index
	Name        string
	Collection  string
	collections Collections
	Pages       map[string]*os.File
}

// NewCollection constracts List of files collection
func NewDatabase(name string) *Database {
	return &Database{
		Name:        rootPath() + name + slash,
		Collection:  "test", // + slash,
		collections: Collections{},
		Pages:       make(map[string]*os.File, 2),
	}
}

// NewCollection constracts List of files collection
func NewCollection(name string) Collection {
	return Collection{
		at:           0,
		primaryIndex: 0,
		cachedIndexs: [][2]int64{},
	}

}

// opnens all collection in Root database folder
func (db *Database) Open() {
	path := db.Name // + db.Collection + slash

	files, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(db.Name, 0744)
		if err != nil {
			eLog.Println("while mkDir", err)
		}
	}

	_, err = os.Stat(db.Name + db.Collection + pi)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(db.Name+db.Collection+pi, os.O_CREATE|os.O_RDWR, 0644)
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

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		page, err := os.OpenFile(path+file.Name(), os.O_RDWR, 0644) //
		iLog.Println("open db path file is : ", path)
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+file.Name(), err)
		}

		db.Pages[db.Name+file.Name()] = page
	}

	if len(db.Pages) < 2 {
		println("path is ", path)
		page, err := os.OpenFile(path+db.Collection+"0", os.O_CREATE|os.O_RDWR, 0644) //
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+"0", err)
		}

		db.Pages[path+db.Collection+"0"] = page
	}
}

// closes All collection
func (db *Database) Close() {
	for _, page := range db.Pages {
		page.Close()
		iLog.Printf("%s closed\n", page.Name())
	}
}
