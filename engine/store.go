package engine

import (
	"encoding/binary"
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

var db *Store

type Store struct {
	db     *bbolt.DB
	lastid map[string]uint64
}

func NewDB(path string) *Store {

	// Open a bbolt database
	kv, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	lastIds := make(map[string]uint64, 0)

	err = kv.View(func(tx *bbolt.Tx) error {
		// Iterate over all buckets in the root
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {

			var lastKey = uint64ToBytes(0)

			// Get the last key in the bucket
			err = kv.View(func(tx *bbolt.Tx) error {
				bucket := tx.Bucket(name)
				if bucket == nil {
					return fmt.Errorf("bucket not found")
				}
				lastKey, _ = bucket.Cursor().Last()
				return nil
			})
			if err != nil {
				return err
			}

			lastIds[string(name)] = binary.BigEndian.Uint64(lastKey)
			return nil
		})

	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	db = &Store{db: kv, lastid: lastIds}
	//toRemove for k, v := range db.lastid {fmt.Println(k, v)}
	return db
}

// inset
func (s *Store) Put(coll, val string) (err error) {
	err = s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(coll))
		if err != nil {
			fmt.Println("err in put db.db.update", err)
			return err
		}
		key := uint64ToBytes(db.lastid[coll])
		bucket.Put(key, []byte(val))
		return nil
	})
	return err
}

func (s *Store) Close() {
	s.db.Sync()
	s.db.Close()
}
