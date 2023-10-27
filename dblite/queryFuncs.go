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
	pindex := db.Name + collection + pIndex

	// if collection == "" {return "ERROR! insert into no collection"}
	_, ok := Indexs[pindex]
	if !ok {
		CreateCollection(collection)
		//return "create " + collection + " first"
	}
	// page name as int
	pName := Indexs[pindex].primaryIndex / MaxObjects

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

	value, err := sjson.Set(data, "_id", Indexs[pindex].primaryIndex)
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
	AppendIndex(db.Pages[pindex], Indexs[pindex].at, size)

	Indexs[pindex].at += int64(size)
	Indexs[pindex].primaryIndex++

	return fmt.Sprint("Success Insert, _id: ", Indexs[pindex].primaryIndex-1)
}

// Select reads data form docs
func SelectById(query string) (result string) {
	// TODO check is collection exist! or make client lib check it
	collection := gjson.Get(query, "collection").String() // + slash

	pindex := db.Name + collection + pIndex
	_, ok := Indexs[pindex]
	if !ok {
		return "Error! " + collection + "is not exists"
		//return "create " + collection + " first"
	}

	id := gjson.Get(query, "where_id").Int()
	// TODO if no where_id in update query then it return 0, it means update obj _id: 0.
	// Solution is initialize primary Index to 1 ensteade 0,
	// Or check lenth of where_id field befor convert it to int
	// or make client lib checkeing this situation

	if int(id) >= len(Indexs[pindex].indexCache) {
		iLog.Println(id, "index not found")
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := Indexs[pindex].indexCache[id][0]
	size := Indexs[pindex].indexCache[id][1]

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

	id := gjson.Get(data, "_id").Int() //fmt.Println("id is : ", id)

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	_, err := Append(db.Pages[path], data)
	if err != nil {
		return fmt.Errorf("ERROR! from Append %v\n", err).Error()
	}

	// Update index
	size := int64(len(data))

	pindex := db.Name + collection + pIndex

	UpdateIndex(db.Pages[pindex], int(id), Indexs[pindex].at, size)

	Indexs[pindex].at += size

	return "Success update"
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

func Select(query string) string {
	return ""
}
