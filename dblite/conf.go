package dblite

import (
	"log"
	"os"
)

var MaxObjects int64 = 10_000

var iLog = log.New(os.Stdout, "\n\033[33mINFO!:  \033[0m", log.Lshortfile)  // log.LstdFlags|
var eLog = log.New(os.Stdout, "\n\033[31mERROR!:  \033[0m", log.Lshortfile) // log.LstdFlags|

type conf struct {
}
