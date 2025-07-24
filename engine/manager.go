package engine

import (
	"github.com/tidwall/gjson"
)

func (s *Store) getCollections() string {
	return "not implemented yet"
}

// deletes collection
func (s *Store) deleteCollection(query gjson.Result) string {
	// TODO return number of deleted objects
	coll := query.Get("collection").Str

	_, err := s.db.Exec("drop table " + coll)
	if err != nil {

		return `{"satatus":"delecte table success"}`
	}
	return "not implemented yet"
}

// creates new collection
func (s *Store) createCollection(query gjson.Result) string {
	coll := query.Get("collection").Str
	s.db.Exec("create table " + coll + "(obj text);")
	s.lastids[coll] = 1
	return coll + " Done"
}

// Rename renames db.
func (s *Store) renameDB(query gjson.Result) string {
	_ = query
	return "renameDB not emplemented yet"
}

// Remove remove db to .Trash dir
func (s *Store) removeDB(query gjson.Result) string {
	_ = query
	return "removeDB not emplemented yet"
}

// ???
func (s *Store) createDB(query gjson.Result) string {
	_ = query
	return "createDB not emplemented yet"
}

func deleteDB(query gjson.Result) string {
	_ = query
	return "removeDB not emplemented yet"
}
