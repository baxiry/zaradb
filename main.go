package main

import (
	"fmt"
)

func main() {
	pages := NewPages()
	pages.Open(RootPath)
	defer pages.Close()

	fmt.Println("pages : ", pages.Pages)
}
