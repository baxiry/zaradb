package main

import (
	"fmt"
)

func main() {

	pages := NewPages()
	fmt.Println("nomber of pages", len(pages.Pages))

	pages.Open(RootPath)
	defer pages.Close()

	fmt.Println("nomber of pages", len(pages.Pages))

	fmt.Println("pages : ", pages.Pages)

	for i := 5; i < 10; i++ {

		SetIndex("123456789 \n", "primary.index")
	}

}
