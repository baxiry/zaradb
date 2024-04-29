package main

import (
	"log"
	"os"
)

const Host = "1111"
const Port = "localhost"

type config struct {
	// TODO
}

const slash = string(os.PathSeparator) // not tested for windos

var iLog = log.New(os.Stdout, "\n\033[33mINFO!:  \033[0m", log.Lshortfile)  // log.LstdFlags|
var eLog = log.New(os.Stdout, "\n\033[31mERROR!:  \033[0m", log.Lshortfile) // log.LstdFlags|
