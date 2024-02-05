package main

import (
	"store"
)

func main() {
	eng := store.NewEngine("mydb")
	for ok := range eng.Collections {

		println(ok)
	}

	eng.Close()
}
