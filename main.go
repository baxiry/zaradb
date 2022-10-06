package main

import (
	"fmt"
)

// TODO delete db ?!

// TODO create collecte
// TODO rename collecte
// TODO delete collecte

// TODO show dbs
// TODO show collects
// TODO switch bitween dbs

func main() {
	dbname, err := CreateDB("testdb1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create db:", dbname)

	colname, err := CreateCl(dbname + "/users1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("create collection :", colname)
}
