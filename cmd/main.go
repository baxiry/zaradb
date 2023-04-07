package main

import (
	"db"
	"fmt"
)

func init() {
	fmt.Println("root path is ", db.RootPath)
}

func main() {

	fmt.Println("Done")
}
