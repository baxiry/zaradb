package database

import (
	"log"
	"os"
)

const slash = string(os.PathSeparator) // not tested for windos

var iLog = log.New(os.Stdout, "\n\033[33mINFO!:  \033[0m", log.Lshortfile)  // log.LstdFlags|
var eLog = log.New(os.Stdout, "\n\033[31mERROR!:  \033[0m", log.Lshortfile) // log.LstdFlags|
