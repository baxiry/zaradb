package dblite

var db *Database

func Run(path string) *Database {
	db = Open(path)
	return db
}
