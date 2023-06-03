package main

import (
	"fmt"
)

func main() {
	fmt.Println("rootPath : ", RootPath)

	pages := NewPages()
	//	fmt.Println(len(pages.Pages))
	pages.Open(RootPath)
	defer pages.Close()

	fmt.Println("pages : ", pages.Pages)

}
