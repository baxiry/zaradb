package main

import (
	"strconv"
	"strings"
)

const LenIndex = 20
const IndexsFile = RootPath + "primary.index"

func StoreIndex(path string) {}

// LastIndex return last index in table
func LastIndex(path string) int {
	last := 0 // read last indext from tail file
	return last + 1
}

// getLocation return fileName & indexLocation where data stored
func GetLocation(id int) (fileName string, indexLocation int64) {
	fileName = strconv.Itoa(id / 1000)
	indexLocation = int64(id % 1000)
	return fileName, indexLocation * LenIndex
}

// convertIndex convert string index location to at and size int64
func ConvIndex(IndexLocation string) (at, size int64) {
	sloc := strings.Split(IndexLocation, " ")
	id, _ := strconv.Atoi(sloc[0])
	siz, _ := strconv.Atoi(sloc[1])
	return int64(id), int64(siz)
}
