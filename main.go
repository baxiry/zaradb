package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func getJson(str string) (json string) {
	var start, end int32

	var i int32
	for i = 0; i < int32(len(str)); i++ {
		if str[i] == '{' {
			start = i
			break
		}
	}
	for i = int32(len(str)) - 1; i >= 0; i-- {
		if str[i] == '}' {
			end = i
			break
		}
	}
	return str[start : end+1]
}

func main() {
	query := arguments()

	switch {
	case len(query) >= 3:
		switch {
		case strings.HasPrefix(query[2], "find"):
			println("arg is find")

		case strings.HasPrefix(query[2], "insert"):

			path := rootPath + query[0] + "/" + query[1] + "/"
			Insert(path, getJson(query[2]))

			println("arg is insert")
			println("data is ", query[2])

		case strings.HasPrefix(query[2], "update"):
			println("arg is update")

		case strings.HasPrefix(query[2], "delete"):
			println("arg is delete")
		} // end switch args[3]

	case len(query) == 3:
		println("Err query not complet")
	case len(query) == 2:

		switch query[1] {

		case "dbs":
			ListDir("")
		case "help":
			println(help_messages)
		default:
			ListDir(query[1])
		}

	default:
		fmt.Println("Finally get default")
	}
}

// documents /////////////////////////////////////////////////

// Update update document data
func Insert(path, data string) (err error) {
	// TODO add ''where'' statment ensteade serial
	serial := GenSerial(LEN_SERIAL)

	fmt.Println("path : ", path+serial)

	err = os.WriteFile(path+serial, []byte(data), 0644)
	if err != nil {
		fmt.Println("Insert ", err)
		return
	}

	f, err := os.OpenFile(path+"/indexes", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	if _, err = f.WriteString(serial + " "); err != nil {
		fmt.Println("Err save indexes ", err)
	}

	return
}

// rootPath = "/Users/fedora/.mydb/test/"

// Select reads data form docs
func Select(path string) (data string, err error) {

	bdata, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bdata), nil
}

// Update update document data
func Update(serial, data string) (err error) {

	// TODO add ''where'' statment ensteade serial

	err = os.WriteFile(serial, []byte(data), 0644)
	if err != nil {
		return
	}
	return
}

// Delete remove document
func Delete(path string) (err error) {
	err = os.Rename(path, ".Crash/"+path)
	if err != nil {
		return err
	}
	return
}

func arguments() (args []string) {
	args = os.Args
	if len(args) < 2 || args[1] == "" {
		fmt.Println("not enought arguments")
		return
	}
	return strings.Split(args[1], ".")
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

// collections //////////////////////////////////////////////////////////////////////
// type Collection string

// CreateCl create collection
func CreateCl(cPath string) (colname string, err error) { // db and collection Path
	err = os.MkdirAll(rootPath+cPath+"/.Trash/", 0755)
	if err != nil {
		return "", err
	}
	f, err := os.Create(rootPath + cPath + "/indexes")
	if err != nil {
		return "", err
	}
	f.Close()
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
