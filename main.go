package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// TODO delete db ?!

// TODO create collecte
// TODO rename collecte
// TODO delete collecte

// TODO show dbs
// TODO show collects
// TODO switch bitween dbs

func main() {
	fmt.Println(CreateDB("universities"))
}

// CreateDB create db TODO return this directly
func CreateDB(dbName string) (err error) {
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {
	return os.MkdirAll(rootPath+dbName+"/.Trash/", 0755)
	// }
	// return err
}

func DeleteDB(dbName string) string {
	return dbName + " db deleted!"
}

// RemoveDB remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

// RenameDB rename db.
func RenameDB(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}

// rootPath = "/Users/fedora/.mydb/test/"

// Update update document data
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
	err = os.Rename(path, ".Crash/"+path)
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
