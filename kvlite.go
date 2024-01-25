package kvlite

import (
	"fmt"

	"github.com/tidwall/wal"
)

type writer struct {
	log   *wal.Log
	index uint64
}

// insert append data
func (w *writer) insert(data string) (uint64, error) {
	err := w.log.Write(w.index+1, []byte(data))
	if err != nil {
		return 0, err
	}
	w.index++
	return w.index, nil
}

func (w *writer) get(i uint64) (string, error) {
	data, err := w.log.Read(i)
	return string(data), err
}

func NewWriter(path string) (*writer, error) {
	opts := &wal.Options{NoSync: true}
	log, err := wal.Open(path, opts)
	if err != nil {
		return nil, err
	}
	index, err := log.LastIndex()
	if err != nil {
		panic(err)
	}
	writer := &writer{
		index: index,
		log:   log,
	}
	return writer, nil
}

// deletes exist value
func (db *writer) Delete(id int, coll string) string {
	return "done"
}

var str = fmt.Sprint

// updates exist value
func (db *writer) Update(id int, coll, value string) string {
	return "done"
}

// rebuilds indexs
// func (db *Database) reIndex() (indexs map[int]index) {
func (db *writer) reIndex() (indexs []uint64) {
	return indexs
}

// save index
func saveIndexs() {

}

// Close db
func (db *writer) Close() {

	saveIndexs()
}
