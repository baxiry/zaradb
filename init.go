package dblite

import (
	"log"
	"os"
)

var iLog = log.New(os.Stdout, "\n\033[33mINFO!:  \033[0m", log.Lshortfile)  // log.LstdFlags|
var eLog = log.New(os.Stdout, "\n\033[31mERROR!:  \033[0m", log.Lshortfile) // log.LstdFlags|
/*
func init() {

	log.SetFlags(0) // Remove the default flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}
*/
