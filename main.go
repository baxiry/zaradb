package main

import (
	"fmt"
)

// this func just for quick test this app
func main() {

	pages := NewPages()
	fmt.Println("nomber of pages", len(pages.Pages))

	pages.Open(RootPath)
	defer pages.Close()

	fmt.Println("nomber of pages", len(pages.Pages))

	fmt.Println("pages : ", pages.Pages)

	path := RootPath + IndexsFile

	for i := 5; i < 10; i++ {

		NewIndex(2333, 1024, pages.Pages[path])
	}

}
