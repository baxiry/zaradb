package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	args := arguments()

	switch {
	case len(args) >= 4:
		switch args[3] {
		case "find":
			println("arg is find")

		case "insert":
			println("arg is insert")

		case "update":
			println("arg is update")

		case "delete":
			println("arg is delete")
		} // end switch args[3]

	case len(args) == 3:
		println("Err query not complet")
	case len(args) == 2:

		switch args[1] {

		case "dbs":
			ListDir("")
		case "help":
			println(help_messages)
		default:
			ListDir(args[1])
		}

	default:
		fmt.Println("Finally get default")
	}
}

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return args
}

func PathExist(subPath string) bool {
	_, err := os.Stat(rootPath + subPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ListDir(path string) {
	dbs, err := os.ReadDir(rootPath + path)
	if err != nil {
		fmt.Println(err)
	}

	dirs := 0
	for _, dir := range dbs {
		if dir.IsDir() && string(dir.Name()[0]) != "." {
			dirs++
			print(dir.Name(), " ")
		}
	}
	if dirs > 0 {
		println()
		return
	}
	println(path, "is impty")
}

// data bases //////////////////////////////////////////////////////

// CreateDB create db TODO return this directly
func CreateDB(dbName string) (dbname string, err error) {
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {return err}

	err = os.MkdirAll(rootPath+dbName+"/.Trash/", 0755)
	if err != nil {
		return dbname, err
	}
	return dbName, nil
}

// DeleteDB delete db (free hard drive).
func DeleteDB(dbName string) string {
	return dbName + " db deleted!"
}

// Remove remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

// Rename rename db.
func RenameDB(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}

// documents /////////////////////////////////////////////////

// rootPath = "/Users/fedora/.mydb/test/"
// Update update document data
func Update(serial, data string) (err error) {

	// TODO add ''where'' statment ensteade serial

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

// collections //////////////////////////////////////////////////////////////////////
// type Collection string

// CreateCl create collection
func CreateCl(cPath string) (colname string, err error) { // db and collection Path
	err = os.MkdirAll(rootPath+cPath+"/.Trash/", 0755)
	return colname, err
}

// TODO Delete delete collection (free hard drive). //
func DeleteCl(cPath string) string {
	return cPath + " collection deleted!"
}

// Remove remove db to .Trash dir
func RemoveCl(cPath string) (err error) {
	return RenameCl(cPath, ".Trash/"+cPath)
}

// Rename rename db.
func RenameCl(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}
