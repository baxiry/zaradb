package main

import (
	"math/rand"
	"sync"
	"time"
)

var help_messages = `command              description:
-----------------------------------------------------------------
bye . . . . . . . . . . exit quit q : exit mydb client.
dbs . . . . . . . . . . show all databases.
use <db_name> . . . . . switch to exist db, or create new db if not exist.
collects  . . . . . . . shw all collections in selected db.`

var (
	wg       sync.WaitGroup
	rootPath = "/Users/fedora/.mydb/"                                                                                                                                                                       //"/Users/fedora/.mydb/test/"
	Latters  = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"} // "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	ListLen  = len(Latters) - 2
)

func init() {
	seedRand()
}

func seedRand() {
	rand.Seed(time.Now().Unix())
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 5)
	//		rand.Seed(time.Now().Unix())
	//	}
	//}()
}
