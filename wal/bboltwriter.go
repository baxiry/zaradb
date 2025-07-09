package main

import (
	"log"

	"go.etcd.io/bbolt"
)

type Store struct {
	db *bbolt.DB
}

func newDB(name string) *Store {

	db, err := bbolt.Open(name, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	store := &Store{db: db}

	return store
}

type ItemData struct {
	Bucket string
	Key    []byte
	Value  []byte
}

type InsertData struct {
	Items []ItemData
	Done  chan error
}

func (s *Store) commitTransaction(insert InsertData) {

	err := s.db.Update(func(tx *bbolt.Tx) error {

		for _, item := range insert.Items {
			b, err := tx.CreateBucketIfNotExists([]byte(item.Bucket))
			if err != nil {
				if insert.Done != nil {
					insert.Done <- err
					close(insert.Done)
				}
				continue
			}

			if err := b.Put(item.Key, item.Value); err != nil {
				if insert.Done != nil {
					insert.Done <- err
					close(insert.Done)
				}
				continue
			}
		}
		return nil
	})

	if err != nil {
		for _, item := range insert.Items {
			_ = item
			select {
			case <-insert.Done:
			default:
				insert.Done <- err
				close(insert.Done)
			}
		}
	} else {
		for _, item := range insert.Items {
			_ = item
			select {
			case <-insert.Done:
			default:
				insert.Done <- nil
				close(insert.Done)
			}
		}
	}
}
