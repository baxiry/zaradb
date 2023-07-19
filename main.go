package main

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
)

const HelpMessage = "dont forget arguments!"
const ErrTypeArg = "agruments must be json!"

// this func just for quick test this app
func main() {

	args := os.Args
	if len(args) == 1 {
		fmt.Println(HelpMessage)
		return
	}

	arg := args[1]

	fmt.Println("arg is ", arg)

	action := gjson.Get(arg, "action")
	data := gjson.Get(arg, "data")
	fmt.Println("action : ", action)
	fmt.Println("data   : ", data)
}

/*
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
*/
