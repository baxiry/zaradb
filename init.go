package main

import (
	"fmt"
	"sync"
)

const rootPath = "/Users/fedora/.mydb/"

func init() {
	fmt.Println(rootPath)
}

var wg sync.WaitGroup

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
