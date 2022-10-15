package main

import (
	"bufio"
	"fmt"
	"log"
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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		args, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		args = strings.Replace(args, "\n", "", 1)

		switch {
		case args == "dbs":
			ListDir("")

		case args == "help":
			println(help_messages)

		case args != "":
			if PathExist(args) {
				ListDir(args)
			}

			println("what ?")

		default:
		}
	}
}

func arguments() string {
	args := os.Args
	if len(args) < 2 {
		return ""
	}
	return args[1]
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
	for _, dir := range dbs {
		if dir.IsDir() {
			print(dir.Name(), " ")
		}
	}
	println()
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

// CreateDB create db TODO return this directly
func CreateCl(cPath string) (colname string, err error) { // db and collection Path
	// _, err = os.Stat("go.mod")
	//	if os.IsNotExist(err) {
	err = os.MkdirAll(cPath+"/.Trash/", 0755)
	return colname, err
}

// Delete delete db (free hard drive).
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
