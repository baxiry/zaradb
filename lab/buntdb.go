package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/tidwall/buntdb"
)

func main() {
	// Open the data.db file. It will be created if it doesn't exist.
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	start := time.Now()
	db.Update(func(tx *buntdb.Tx) error {
		for i := 0; i < 1000000; i++ {

			tx.Set(strconv.Itoa(i), `{"name":{"first":"Tom","last":"Joh\nnson"},"age":38}`, nil)
		}
		return nil
	})
	fmt.Println("write duration: ", time.Since(start))

	start = time.Now()
	db.View(func(tx *buntdb.Tx) error {
		var data string
		var err error
		var ln int

		for i := 0; i < 1000000; i++ {
			data, err = tx.Get("1")
			ln += len(data)

		}
		fmt.Println("len data is ", ln)
		return err
	})

	fmt.Println("read duration : ", time.Since(start))
}
