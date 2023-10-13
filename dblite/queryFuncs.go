package dblite

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// nuber page to make namePage from id
var numberPage int64 = 0

const slash = string(os.PathSeparator) // not tested for windos

// Insert
func Insert(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash

	// if collection == "" {return "ERROR! insert into no collection"}

	pName := indexs["test"].primaryIndex / MaxObjects // page name as int

	if pName != numberPage {
		numberPage++

		pagePath := db.Name + collection + fmt.Sprint(pName)

		page, err := os.OpenFile(pagePath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("os open file: ", err)
		}

		db.Pages[pagePath] = page
	}

	data := gjson.Get(query, "data").String()

	value, err := sjson.Set(data, "_id", indexs["test"].primaryIndex)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	path := db.Name + collection + fmt.Sprint(pName)

	size, err := Append(db.Pages[path], value)
	if err != nil {
		// TODO check if collection exist
		eLog.Printf("%v Path is %s ", err, path)
		fmt.Println("collection is ", collection)
		return "Fielure Insert,mybe collection not exist"
	}

	// set new index
	AppendIndex(db.Pages[db.Name+collection+pIndex], indexs["test"].at, size)

	indexs["test"].at += int64(size)
	indexs["test"].primaryIndex++
	return fmt.Sprint("Success Insert, _id: ", indexs["test"].primaryIndex-1)
}

// Select reads data form docs
func SelectById(query string) (result string) {
	id := gjson.Get(query, "where_id").Int()
	if int(id) >= len(indexs["test"].indexCache) {
		iLog.Println(id, "index not found")
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := indexs["test"].indexCache[id][0]
	size := indexs["test"].indexCache[id][1]

	collection := gjson.Get(query, "collection").String() // + slash
	//fmt.Println("table is : ", in)
	// TODO check is from exist!

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	return Get(db.Pages[path], at, int(size))
}

// delete
func DeleteById(query string) (result string) {

	id := gjson.Get(query, "_id").Int()
	in := gjson.Get(query, "collection").String() // + slash

	path := db.Name + in + fmt.Sprint(indexs["test"].primaryIndex/MaxObjects)

	fmt.Println("path id DeleteById: ", path)

	UpdateIndex(db.Pages[path], int(id), 0, 0)

	//fmt.Println(IndexsCache.indexs)
	indexs["test"].indexCache[id] = [2]int64{0, 0}
	//fmt.Println(IndexsCache.indexs)

	return "Delete Success!"
}

// Update update document data
func Update(query string) (result string) {
	collection := gjson.Get(query, "collection").String() // + slash
	if collection == "" {
		return "ERROR! select no collection "
	}

	data := SelectById(query)
	newData := gjson.Get(query, "data").String()

	data = gjson.Get("["+data+","+newData+"]", "@join").String()

	id := gjson.Get(data, "_id").Int()

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	_, err := Append(db.Pages[path], data)
	if err != nil {
		return fmt.Errorf("ERROR! from Append %v\n", err).Error()
	}

	// Update index
	size := int64(len(data))

	UpdateIndex(db.Pages[db.Name+collection+pIndex], int(id), indexs["test"].at, size)

	indexs["test"].at += size

	return "Success update"
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

func Select(query string) string {
	return ""
}
