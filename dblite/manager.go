package dblite

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func DeleteCollection(query string) string {
	collectName := gjson.Get(query, "collection").String()

	if collectName == "" {
		return "please type a collection you want to delete"
	}

	// read cllections file name
	infos, err := os.ReadFile(db.Name + db.Infos) //OpenFile(, os.O_RDONLY, 0644)
	if err != nil {
		eLog.Println("open collectFiles : ", err)
	}

	// remove all
	i := 0
	for range db.Pages {
		file := collectName + strconv.Itoa(i)
		i++
		for path := range db.Pages {
			if path == db.Name+file {
				_ = os.Remove(path)
				delete(db.Pages, path)
			}
		}
	}

	_ = os.Remove(db.Name + collectName + pi)
	delete(db.Pages, db.Name+collectName+pi)
	fmt.Println(db.Pages)

	// remove collection from  infos

	cols := strings.TrimRight(string(infos), " ")

	collectsList := strings.Split(cols, " ")

	res := ""
	for _, coll := range collectsList {
		if coll == collectName {
			continue
		}
		res += coll + " "
	}
	fmt.Println("collects is :", res)

	// Write the new content to the file.
	err = os.WriteFile(db.Name+db.Infos, []byte(res), 0644)
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range db.Pages {
		fmt.Println(k, "  : ", v)
	}

	return collectName + " is deleted"
}

// creates new collection
func CreateCollection(query string) string {
	collectName := gjson.Get(query, "collection").String()

	if collectName == "" {
		return "please shose a collectin name"
	}

	// read cllections file name
	infos, err := os.OpenFile(db.Name+db.Infos, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		eLog.Println("open collectFiles : ", err)
	}
	defer infos.Close()

	names, err := io.ReadAll(infos)
	if err != nil {
		eLog.Println(err)
	}

	collectsList := strings.Split(string(names), " ")

	for _, coll := range collectsList {
		if coll == collectName {
			return collectName + " already exist!"
		}
	}

	// create index of collection & first page.

	firstPage, err := os.OpenFile(db.Name+collectName+fmt.Sprint(0), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}
	db.Pages[db.Name+collectName+fmt.Sprint(0)] = firstPage

	indxPage, err := os.OpenFile(db.Name+collectName+pi, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err.Error()
	}
	db.Pages[db.Name+collectName+pi] = indxPage

	_, err = infos.WriteString(collectName + " ")
	if err != nil {
		return fmt.Sprintf("ERROR can't create %s collectiln", collectName)
	}

	for k, v := range db.Pages {
		fmt.Println(k, "  : ", v)
	}

	return "collecteon " + collectName + " is created"
}
