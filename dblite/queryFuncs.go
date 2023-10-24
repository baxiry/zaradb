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

	collection = gjson.Get(query, "collection").String() // + slash

	// if collection == "" {return "ERROR! insert into no collection"}
	_, ok := Indexs[collection+pIndex]
	if !ok {
		CreateCollection(collection)
		//return "create " + collection + " first"
	}
	// page name as int
	pName := Indexs[collection+pIndex].primaryIndex / MaxObjects

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
	if data == "" {
		return "there is no data to insert"
	}

	value, err := sjson.Set(data, "_id", Indexs[collection+pIndex].primaryIndex)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	path := db.Name + collection + fmt.Sprint(pName)

	size, err := Append(db.Pages[path], value)
	if err != nil {
		// TODO check if collection exist
		eLog.Printf("%v\n Path is %s\n collection is %s\n", err, path, collection)
		return "Fielure Insert,mybe collection not exist"
	}

	for k, v := range Indexs {
		fmt.Printf("at in %s is %d\n", k, v.at)
	}
	// set new index
	AppendIndex(db.Pages[db.Name+collection+pIndex], Indexs[collection+pIndex].at, size)

	Indexs[collection+pIndex].at += int64(size)
	Indexs[collection+pIndex].primaryIndex++

	iLog.Printf("coll: %s, primaryIndex: %d\n", collection, Indexs[collection+pIndex].primaryIndex)

	return fmt.Sprint("Success Insert, _id: ", Indexs[collection+pIndex].primaryIndex-1)
}

// Select reads data form docs
func SelectById(query string) (result string) {
	collection := gjson.Get(query, "collection").String() // + slash

	id := gjson.Get(query, "where_id").Int()

	if int(id) >= len(Indexs[collection+pIndex].indexCache) {
		iLog.Println(id, "index not found")
		return fmt.Sprintf("Not Found _id %v\n", id)
	}
	//iLog.Println("collection ; ", collection)

	at := Indexs[collection+pIndex].indexCache[id][0]
	size := Indexs[collection+pIndex].indexCache[id][1]

	//fmt.Println("table is : ", in)
	// TODO check is from exist!

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	return Get(db.Pages[path], at, int(size))
}

// delete
func DeleteById(query string) (result string) {

	collection = gjson.Get(query, "collection").String() // + slash
	// check collection

	id := gjson.Get(query, "where_id").Int()
	fmt.Println("id is : ", id)

	UpdateIndex(db.Pages[db.Name+collection+pIndex], int(id), 0, 0)

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
	fmt.Println("id is : ", id)

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	_, err := Append(db.Pages[path], data)
	if err != nil {
		return fmt.Errorf("ERROR! from Append %v\n", err).Error()
	}

	// Update index
	size := int64(len(data))

	UpdateIndex(db.Pages[db.Name+collection+pIndex], int(id), Indexs[collection].at, size)

	Indexs[collection+pIndex].at += size

	return "Success update"
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

func Select(query string) string {
	return ""
}
