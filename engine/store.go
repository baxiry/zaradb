package engine

import (
	"fmt"
	"log"
	"strconv"

	"go.etcd.io/bbolt"
)

type Store struct {
	db     *bbolt.DB
	lastid map[string]int64
}

var db *Store

func NewDB(path string) *Store {
	// Open a bbolt database
	kv, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		fmt.Println("can not create file ?? why ?")
		log.Fatal(err)
	}

	db = &Store{db: kv, lastid: make(map[string]int64, 0)}
	return db
}

func (db *Store) getLastKey(bucket string) int64 {
	var lastKey = []byte("0")

	// Get the last key in the bucket
	db.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket != nil {
			lastKey, _ = bucket.Cursor().Last()
		}
		return nil
	})
	id, _ := strconv.Atoi(string(lastKey))
	return int64(id)
}

func (db *Store) Put(backet, key, val string) (err error) {
	fmt.Println("key val Is : ", key, val)
	db.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(backet))
		if err != nil {
			fmt.Println("err in Put", err)
			return err
		}
		bucket.Put([]byte(key), []byte(val))
		return nil
	})
	return err
}

func (db *Store) Close() {
	db.db.Sync()
	db.db.Close()
}
