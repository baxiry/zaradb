package dblite

import (
	"os"
)

// a collection fot test
const testCollection = "test"

var collection string //= "test"

// db
type Database struct {
	// name of database
	Name string

	// file of collections name
	collection string

	// slice of collection's name in db
	CollectsList []string

	Pages map[string]*os.File
}

// NewDatabase create new *database
func NewDatabase(name string) *Database {
	return &Database{
		Name:         rootPath() + name + slash,
		CollectsList: make([]string, 0),
		Pages:        make(map[string]*os.File, 2),
	}
}

// opnens all collection in Root database folder
func (db *Database) Open() {
	path := db.Name // + db.Index + slash

	files, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(db.Name, 0744)
		if err != nil {
			eLog.Println("while mkDir", err)
		}
	}

	_, err = os.Stat(db.Name + testCollection + pIndex)
	if os.IsNotExist(err) {
		f, err := os.OpenFile(db.Name+testCollection+pIndex, os.O_CREATE|os.O_RDWR, 0644)
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
		if file.IsDir() || file.Name() == "infos" {
			continue
		}

		page, err := os.OpenFile(path+file.Name(), os.O_RDWR, 0644) //
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+file.Name(), err)
		}

		db.Pages[db.Name+file.Name()] = page
	}

	if len(db.Pages) < 2 {
		println("path is ", path)
		page, err := os.OpenFile(path+testCollection+"0", os.O_CREATE|os.O_RDWR, 0644) //
		if err != nil {
			iLog.Printf("os open file: %s,  %v\n", path+"0", err)
		}

		db.Pages[path+testCollection+"0"] = page
	}
}

// closes All collection
func (db *Database) Close() {
	for _, page := range db.Pages {
		page.Close()
		iLog.Printf("%s closed\n", page.Name())
	}

	// TODO Close all collections in database
	// TODO delete indexs cache
}

// end
