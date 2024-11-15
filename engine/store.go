package engine

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/tidwall/gjson"
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

// getData fitch for data
func (s *Store) getData(query gjson.Result) (data []string, err error) {
	coll := query.Get("collection").Str
	if coll == "" {
		return nil, fmt.Errorf(`{"error":"forgot collection name "}`)
	}

	skip := query.Get("skip").Int()
	limit := query.Get("limit").Int()
	if limit == 0 {
		// what default setting should be here ?
		limit = 1000
	}

	err = s.db.View(func(tx *bbolt.Tx) error {

		bucket := tx.Bucket([]byte(coll))
		if bucket == nil {
			return fmt.Errorf("collection %s is not exists", coll)
		}
		isMatch := query.Get("match")
		// Use a cursor to iterate over all key-value pairs in the bucket.
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {

			if limit == 0 {
				break
			}

			ok, err := match(isMatch, string(value))
			if err != nil {
				return err
			}

			if ok {
				if skip != 0 {
					skip--
					continue
				}
				data = append(data, string(value))
				limit--
			}

		}

		return nil
	})

	return data, err
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

// to work with writer
var (
	oksChan = make(chan []bool, 1)
	done    = make(chan bool, 1)
	objChan = make(chan Object, 1)
)

type Object struct {
	id       int
	bucket   string
	key, val []byte
}

// Function to gather data and send after a duration
func writer(db *bbolt.DB, input chan Object) {
	dataBatch := make(map[string][]Object)
	ticker := time.NewTicker(100 * time.Millisecond)
	var obj Object
	oks := make([]bool, 0)
	for {
		select {
		case obj = <-input:
			dataBatch[obj.bucket] = append(dataBatch[obj.bucket], obj)
			oks = append(oks, true)

		case <-ticker.C:
			if len(dataBatch) == 0 {
				fmt.Println("nothing to write")
				continue
			}
			err := db.Batch(func(tx *bbolt.Tx) error {
				for bucketName, keyValues := range dataBatch {
					bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
					if err != nil {
						return err
					}
					for _, v := range keyValues {
						if err := bucket.Put(v.key, v.val); err != nil {
							return err
						}
					}
				}
				return nil
			})
			if err != nil {
				fmt.Println("error at db.batch", err)
			}
			oksChan <- oks
			dataBatch = map[string][]Object{}
			oks = []bool{}
		case <-done:
			break
		}
	}
}

func toByte(number int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(number))
	return buf
}

func (s *Store) Close() {
	s.db.Sync()
	s.db.Close()
}
