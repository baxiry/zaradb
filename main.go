package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func main() {
	path := "example.db"
	file, err := Opendb(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	src := ""

	for i := 0; i < 1000; i++ {
		AppendData(file, genData(i))
	}

	fmt.Println("size", FileSize(path))

	for i := 0; i < 1000; i++ {
		src = getVal(file, int64(10*i), 10)
		fmt.Println(src)
	}

}

func getId(json string) string {
	value := gjson.Get(json, "_id")
	println(value.String())
	return value.String()
}

func genData(n int) (data string) {
	num := strconv.Itoa(n)
	data = num
	for i := 0; i < 10-len(num); i++ {
		data += "_"

	}
	return data
}

// AppendData to file
func AppendData(file *os.File, data string) (err error) {
	lnb, err := file.WriteString(data)
	println("len bytes is : ", lnb)
	return
}

func getVal(file *os.File, at int64, ln int) string {
	buffer := make([]byte, ln)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println("file size is : ", FileSize(file.Name()))
		fmt.Println("at is ", at)
		panic(err)
	}
	// out the buffer content
	return string(buffer[:n])
}

func Opendb(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

func FileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return file.Size()
}

func IsExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
func queryLang() {
	query := arguments()
	fmt.Println("query is : ", query)

	switch {
	case len(query) >= 3:
		switch {
		case strings.HasPrefix(query[2], "find"):
			fmt.Println(" find by ", getJson(query[2]))
			getId(query[2])

			d, err := Select(rootPath + query[0] + "/" + query[1] + "/" + getId(query[2]))
			if err != nil {
				fmt.Println("error is :", err)
			}

			fmt.Println("data is ", d)

		case strings.HasPrefix(query[2], "insert"):
			path := rootPath + query[0] + "/" + query[1] + "/"
			Insert(path, getJson(query[2]))
			println("insert", getJson(query[2]))

		case strings.HasPrefix(query[2], "update"):
			println("arg is update")

		case strings.HasPrefix(query[2], "delete"):
			println("arg is delete")
		} // end switch args[3]

	case len(query) == 2:
		println("Err query not complet")
	case len(query) == 1:

		switch query[0] {

		case "dbs":
			ListDir("")
		case "help":
			println(help_messages)
		default:
			ListDir(query[0])
		}

	default:
		fmt.Println("Finally get default")
	}
}

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

// documents /////////////////////////////////////////////////

// Update update document data
func Insert(path, data string) (err error) {
	// TODO add ''where'' statment insteade by serial
	return
}

// rootPath = "/Users/fedora/.mydb/test/"

// Select reads data form docs
func Select(path string) (data string, err error) {
	return data, nil
}

// Update update document data
func Update(serial, data string) (err error) {
	// TODO add ''where'' statment ensteade serial
	return
}

// Delete remove document
func Delete(path string) (err error) {
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

func getIndexes(path string) []string {
	return []string{}
}

// Rename rename db.
func RenameDB(oldPath, newPath string) (err error) {
	return os.Rename(oldPath, newPath)
}

// Remove remove db to .Trash dir
func RemoveDB(dbName string) (err error) {
	return RenameDB(dbName, ".Trash/"+dbName)
}

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
