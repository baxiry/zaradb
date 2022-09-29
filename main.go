package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
)

// TODO create db
// TODO rename db
// TODO delete db

// TODO create collecte
// TODO rename collecte
// TODO delete collecte

// TODO show dbs
// TODO show collects
// TODO switch bitween dbs

func CreateDB(dbName string) (err error) {
	_, err = os.Stat(rootPath + dbName)

	if os.IsNotExist(err) {
		err = os.MkdirAll(rootPath+dbName, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil

}
func main() {
	err := CreateDB("dbdir")
	fmt.Println(err)
}

// rootPath = "/Users/fedora/.mydb/test/"
// Update data
func Update(serial, data string) (err error) {

	// TODO add ''where'' statment enstead serial

	err = os.WriteFile(serial, []byte(data), 0644)
	if err != nil {
		return
	}
	return
}

// Select reads data form docs
func Select(path string) (data string, err error) {

	bdata, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bdata), nil

}

// Delete remove document
func Delete(path string) (err error) {
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return
}

// GenSerial generate serial for Doc
func GenSerial(length int) (serial string) {
	var i int
	for i = 0; i < length; i++ {
		serial += Latters[rand.Intn(ListLen)+1]
	}
	return serial
}

func getIndexes(path string) []string {
	data, err := Select(path)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(string(data), " ")
}
