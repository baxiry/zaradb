package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func main() {
	file, err := Opendbs("example.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	start := time.Now()
	data := ""
	size := 0
	for i := 0; i < 10000000; i++ {
		data = " hello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello wordhello word"
		size += len(data)
		AppendData(file, strconv.Itoa(i)+data)
	}

	fmt.Println("write duration: ", time.Since(start))

	start = time.Now()
	leen := 0
	op := 0

	for i := 0; i < 10000000000; i += 100 {
		leen += len(GetVal(file, int64(i), len(data)))
		op++
	}

	fmt.Println("len of geted data : ", leen)
	fmt.Println("no op", op)
	fmt.Println("duration read : ", time.Since(start))
}

func openIndexes() {}

// Opendb opens | create  new file db
func Opendbs(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	return file, err
}

// GetVal return data string.
// take file pointr, at int64 & len of data will read
func GetVal(file *os.File, at int64, buff int) string {
	buffer := make([]byte, buff)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return ""
	}
	// out the buffer content
	return string(buffer[:n])
}

// AppendData appends data string to file
// return len or size of file and err
func AppendData(file *os.File, data string) (int, error) {
	lenByte, err := file.WriteString(data)
	return lenByte, err
}

// getField get json field from json string
func getField(field, json string) string {
	value := gjson.Get(json, field)
	//println(value.String())
	return value.String()
}

// check if path is exist
func IsExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// simplest query language
func queryLang() {
	query := arguments()
	fmt.Println("query is : ", query)
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

// /////////////////////////////////////////////////////////////////////////////////////
func PathExist(subPath string) bool {
	_, err := os.Stat(rootPath + subPath)
	if os.IsNotExist(err) {
		return false
	}
	return true
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
