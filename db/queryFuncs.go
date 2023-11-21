package db

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// nuber page to make namePage from id
var numberPage int64 = 0

const slash = string(os.PathSeparator) // not tested for windos

// Find finds any obs match creteria.
func findMany(query string) (res string) {

	collection := gjson.Get(query, "collection").String()
	_ = collection

	filter := gjson.Get(query, "filter").String()

	limit := int64(20)
	// offset := 0

	// reads first 20 item by default

	listObj := make([]string, limit)

	var i int64
	var ii int64

	for i = 0; ii < limit; i++ {
		data := "get rose"
		if match(filter, data) {
			listObj[ii] = data // + ",\n"
			ii++
		}
	}

	res = "[\n"
	for k, v := range listObj {
		if v == "" {
			fmt.Println("zero val")
			break
		}
		res += " " + listObj[k] + ",\n"
	}
	if len(res) == 2 {
		return "[]"
	}
	return res[:len(res)-2] + "\n]"
}

// Finds first obj match creteria.
func findOne(query string) (res string) {
	collection := gjson.Get(query, "collection").String()

	filter := gjson.Get(query, "filter").String()

	var i int64

	for i = 0; i < 10; i++ {

		if match(filter, res) {
			return res
		}
	}

	return "now data match" + collection

}

// findById reads data form docs
func findById(query string) string {
	return "data"
}

// Insert
func Insert(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash

	data := gjson.Get(query, "data").String()

	value, err := sjson.Set(data, "_id", "key")
	if err != nil {
		fmt.Println("sjson.Set : ", err)
	}

	return fmt.Sprint("Success Insert, _id: ", collection, value)
}

// delete
func DeleteById(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash
	// check collection

	id := gjson.Get(query, "where_id").Int()
	fmt.Println("id is : ", id)

	return "Delete Success!" + collection
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
	_ = id

	return "Success update"
}

// end
