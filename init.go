package main

import (
	"math/rand"
	"sync"
	"time"
)

var help_messages = `command & description:
-----------------------------------------------------------------
help :
   get this help message.

dbs  :
   show exist databases.

<db_name> : 
   shw all collections in selected db. Note: collection is table.

dbName.collectionName.find :
   find * record in dbName.collectionName .

dbName.collectionName.insert {...}:
   insert new record to dbName.collection.
   Note: make sure that data is a json. 
`

const (
	LEN_SERIAL = 6
)

var (
	wg       sync.WaitGroup
	rootPath = "/Users/fedora/.mydb/"
	Latters  = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	ListLen  = len(Latters) - 2
)

func init() {
	seedRand()
}

func seedRand() {
	rand.Seed(time.Now().Unix())
}
