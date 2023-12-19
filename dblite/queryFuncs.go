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

// Finds first obj match creteria.
func findOne(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash

	for i := 0; i <= db.Lid; i++ {
		if db.indexs[i].coll != collection {
			continue
		}
		data := db.Get(i, collection)
		filter := gjson.Get(query, "filter").String()
		if match(filter, data) {
			return data
		}
	}

	return "noting mutch"
}

// Find finds any obs match creteria.
func findMany(query string) (res string) {
	return "not emplement yet"
}

// findById reads data form docs
func findById(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash

	id := gjson.Get(query, "_id").Int()

	return db.Get(int(id), collection)
}

// Insert
func Insert(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash
	//	CreateCollection(collection)

	data := gjson.Get(query, "data").String()
	if data == "" {
		return "there is no data to insert"
	}

	value, err := sjson.Set(data, "_id", db.Lid+1)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	db.Insert(collection, value)

	return fmt.Sprint("Success Insert, _id: ", db.Lid)
}

// delete
func DeleteById(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash
	// check collection

	id := gjson.Get(query, "where_id").Int()

	db.Delete(int(id), collection)

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

	db.Update(int(id), collection, data)

	return "Success update"
}

// end
