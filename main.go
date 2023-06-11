package main

import (
	"fmt"
)

func main() {
	fmt.Println("rootPath : ", RootPath)

	pages := NewPages()
	fmt.Println("nomber of pages", len(pages.Pages))
	pages.Open(RootPath)
	defer pages.Close()

	fmt.Println("nomber of pages", len(pages.Pages))
	fmt.Println("pages : ", pages.Pages)

}
