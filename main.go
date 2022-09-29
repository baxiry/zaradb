package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// TODO rename db
// TODO delete db

// TODO create collecte
// TODO rename collecte
// TODO delete collecte

// TODO show dbs
// TODO show collects
// TODO switch bitween dbs

func RenameDB(oldPath, newPath string) string {

	err := os.Rename(oldPath, newPath)

	if err != nil {
		return err.Error()
	}
	return "Done"
}

func main() {

	println(mydb_art)
	println(mydb_art1)
	println(mydb_art2)
	err := CreateDB("peoples")
	fmt.Println(err)
	//	err = RenameDB("peoples", "books")
	fmt.Println(err)
}

func DeleteDB(path string) string {

	return ""
}
func CreateDB(dbName string) string {
	_, err := os.Stat(rootPath + dbName)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rootPath+dbName+"/.Crash/", 0755)
		if err != nil {
			return err.Error()
		}
	}
	return "Done"
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
