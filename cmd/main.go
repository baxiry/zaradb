package main

import "store"

func main() {
	eng := store.NewDatabase("mydb")
	for ok := range eng.Collections {

		println(ok)
	}

	eng.Close()
}
