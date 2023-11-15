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

func findOne(query string) (res string) {
	collection = gjson.Get(query, "collection").String()

	filter := gjson.Get(query, "filter").String()

	pindex := db.Name + collection + pIndex

	path := ""

	var i int64
	for i = 0; i < int64(len(Indexs)); i++ {

		at := Indexs[pindex].indexCache[i][0]
		size := Indexs[pindex].indexCache[i][1]

		// TODO check performence of this
		path = db.Name + collection + fmt.Sprint(i/MaxObjects)

		res = Get(db.Pages[path], at, int(size))
		if match(filter, res) {
			fmt.Println("res:  ", res)
			fmt.Println("\nfilter", filter)
			return res
		}
	}

	return "now data match"

}

// Find finds many by filter.
func findMany(query string) (res string) {
	// if sub index not exists

	coll := gjson.Get(query, "collection").String()
	// if len(coll) == 0 {return "select collection"}

	filter := gjson.Get(query, "filter").String()
	fmt.Println("filter is :", filter)

	pindex := db.Name + coll + pIndex

	limit := int64(20)
	// offset := 0
	if int(limit) >= len(Indexs[pindex].indexCache) /* -offset */ {
		limit = int64(len(Indexs[pindex].indexCache)) - 1
	}

	// reads first 20 item by default

	listObj := make([]string, limit)

	var i int64
	for i = 0; i < limit; i++ {

		at := Indexs[pindex].indexCache[i][0]
		size := Indexs[pindex].indexCache[i][1]

		path := db.Name + coll + fmt.Sprint(i/MaxObjects)

		listObj[i] = Get(db.Pages[path], at, int(size)) // + ",\n"
	}

	res = "[\n"
	for i := 0; i < int(limit); i++ {
		res += listObj[i] + ",\n"
	}

	return res[:len(res)-2] + "\n]"
}

// findById reads data form docs
func findById(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash

	pindex := db.Name + collection + pIndex

	if _, ok := Indexs[pindex]; !ok {
		// TODO union all expected errors
		return "Error! " + collection + " is not exists"
		//return "create " + collection + " first"
	}

	id := gjson.Get(query, "where_id").Int()

	if int(id) >= len(Indexs[pindex].indexCache) {
		iLog.Println(id, "index not found")
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := Indexs[pindex].indexCache[id][0]
	size := Indexs[pindex].indexCache[id][1]
	if size == 0 {
		return ""
	}

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	return Get(db.Pages[path], at, int(size))
}

// Insert
func Insert(query string) (res string) {

	collection = gjson.Get(query, "collection").String() // + slash
	pindex := db.Name + collection + pIndex

	// if collection == "" {return "ERROR! insert into no collection"}
	_, ok := Indexs[pindex]
	if !ok {
		//return "create " + collection + " first"
		CreateCollection(collection)
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
		return "Fielure Insert, mybe collection is not exist"
	}

	// store the new index
	AppendIndex(db.Pages[pindex], Indexs[pindex].at, size)

	Indexs[pindex].at += int64(size)
	Indexs[pindex].primaryIndex++

	return fmt.Sprint("Success Insert, _id: ", Indexs[pindex].primaryIndex-1)
}

// delete
func DeleteById(query string) string {

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

	data := findById(query)
	newData := gjson.Get(query, "data").String()

	data = gjson.Get("["+data+","+newData+"]", "@join").String()

	id := gjson.Get(data, "_id").Int()
	// TODO if no where_id in update query then it return 0, it means update obj _id: 0.
	// Solution is initialize primary Index to 1 insteade 0,
	// Or check length of where_id field befor convert it to int
	// or make client lib checkeing this situation

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

// end
