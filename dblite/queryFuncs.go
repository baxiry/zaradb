package dblite

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var MaxObjects int64 = 10_000

// nuber page to make namePage from id
var numberPage int64 = 0

const slash = string(os.PathSeparator) // not tested for windos

// Insert
func Insert(query string) (res string) {

	collection := gjson.Get(query, "collection").String() + slash
	if len(collection) == 1 {
		return fmt.Sprint("failure insert. insert into no collection")
	}

	pName := db.PrimaryIndex / MaxObjects // page name as int

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

	value, err := sjson.Set(data, "_id", db.PrimaryIndex)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	path := db.Name + collection + fmt.Sprint(pName)

	size, err := Append(db.Pages[path], value)
	if err != nil {
		// TODO check if collection exist
		eLog.Printf("%v Path is %s ", err, path)
		return "Fielure Insert,mybe collection not exist"
	}

	// set new index
	AppendIndex(db.Pages[db.Name+collection+pi], collect.at, size)

	collect.at += int64(size)
	db.PrimaryIndex++
	return fmt.Sprint("Success Insert, _id: ", db.PrimaryIndex-1)
}

// Select reads data form docs
func SelectById(query string) (result string) {
	id := gjson.Get(query, "where_id").Int()
	if int(id) >= len(collect.cachedIndexs) {
		iLog.Println("no found index", id)
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := collect.cachedIndexs[id][0]
	size := collect.cachedIndexs[id][1]

	in := gjson.Get(query, "collection").String() + slash
	//fmt.Println("table is : ", in)
	// TODO check is from exist!

	path := db.Name + in + fmt.Sprintf("%d", id/MaxObjects)

	return Get(db.Pages[path], at, int(size))
}

// delete
func DeleteById(query string) (result string) {

	id := gjson.Get(query, "_id").Int()
	in := gjson.Get(query, "collection").String() + slash

	path := db.Name + in + fmt.Sprint(db.PrimaryIndex/MaxObjects)

	fmt.Println("path id DeleteById: ", path)

	UpdateIndex(db.Pages[path], int(id), 0, 0)

	//fmt.Println(IndexsCache.indexs)
	collect.cachedIndexs[id] = [2]int64{0, 0}
	//fmt.Println(IndexsCache.indexs)

	return "Delete Success!"
}

// Update update document data
func Update(query string) (result string) {
	collection := gjson.Get(query, "collection").String() + slash
	if len(collection) == 1 {
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

	UpdateIndex(db.Pages[db.Name+collection+pi], int(id), collect.at, size)

	collect.at += size

	return "Success update"
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

func selectFields(query string) string {
	return ""
}
