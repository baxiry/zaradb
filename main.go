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
	fmt.Println(CreateDB("testdb"))
	fmt.Println(Create("users"))
}
