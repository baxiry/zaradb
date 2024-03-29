package db

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Finds first obj match creteria.
func findOne(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash

	for i := 0; i <= db.lastId; i++ {
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

	collection := gjson.Get(query, "collection").String() // + slash

	limit := int(gjson.Get(query, "limit").Int())
	if limit == 0 {
		limit = 100
	} else if limit > db.lastId {
		limit = 200
	}

	skip := int(gjson.Get(query, "skip").Int())
	if skip > db.lastId {
		skip = db.lastId - 1
		println("skip greater then last id", skip)
	}

	res = "["

	for i := skip + 1; i <= db.lastId; i++ {
		if db.indexs[i].coll != collection {
			continue
		}
		data := db.Get(i, collection)
		filter := gjson.Get(query, "filter").String()
		if match(filter, data) {
			res += data + ","
			limit--
			if limit == 0 {
				break
			}
		}
	}
	if len(res) == 1 {
		return "[]"
	}
	res = res[:len(res)-1] + "]"
	return res
}

// findById reads data form docs
func findById(query string) string {
	for i := 0; i <= db.lastId; i++ {
		fmt.Println(i, db.indexs[i])
	}
	collection := gjson.Get(query, "collection").String() // + slash

	id := gjson.Get(query, "_id").Int()

	return db.Get(int(id), collection)
}

// Insert
func insert(query string) (res string) {

	collection := gjson.Get(query, "collection").String() // + slash
	//	CreateCollection(collection)

	data := gjson.Get(query, "data").String()
	if data == "" {
		return "there is no data to insert"
	}

	value, err := sjson.Set(data, "_id", db.lastId+1)
	if err != nil {
		fmt.Println("sjson.Set : ", err)
		return "internal error"
	}

	// make this return error
	db.Insert(collection, value)

	return fmt.Sprint("Success Insert, _id: ", db.lastId)
}

// delete
func deleteOne(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash
	// check collection
	for i := 0; i <= db.lastId; i++ {
		if db.indexs[i].size == 0 {
			println("skiping id: ", i)
			continue
		}

		// Mach
		filter := gjson.Get(query, "filter").String()
		data := db.Get(i, collection)
		if match(filter, data) {
			println("delete id: ", i)
			return db.Delete(i, collection)
		}
	}
	return "nothing match"
}

// delete
func deleteMany(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash
	// check collection

	// indx, ok := db.indexs[id]; if !ok { return "no data to delete"	}
	tot := 0

	for i := 0; i <= db.lastId; i++ {
		if db.indexs[i].size == 0 {
			continue
		}

		if db.indexs[i].coll != collection {
			continue
		}

		// Mach
		filter := gjson.Get(query, "filter").String()
		data := db.Get(i, collection)
		if match(filter, data) {

			if db.Delete(i, collection) == "delete success!" {
				tot++
			}
		}
	}
	return str(tot) + " items deleted!"
}

// delete by id
func deleteById(query string) string {

	collection := gjson.Get(query, "collection").String() // + slash

	id := gjson.Get(query, "_id").Int()

	return db.Delete(int(id), collection)
}

// Update update document data
func updateById(query string) (result string) {
	collection := gjson.Get(query, "collection").String() // + slash
	if collection == "" {
		return "ERROR! select no collection "
	}

	// TODO make findById return error
	data := findById(query)

	newData := gjson.Get(query, "data").String()
	data = gjson.Get("["+data+","+newData+"]", "@join").String()

	id := gjson.Get(data, "_id").Int()

	db.Update(int(id), collection, data)

	return "Success update"
}

// TODO updateOne one update document data
func updateOne(query string) (result string) {
	collection := gjson.Get(query, "collection").String() // + slash
	if collection == "" {
		return "ERROR! select no collection "
	}

	filter := gjson.Get(query, "filter").String()
	newData := gjson.Get(query, "data").String()

	for i := 0; i <= db.lastId; i++ {

		data := db.Get(i, collection)

		if match(filter, data) {

			data = gjson.Get("["+data+","+newData+"]", "@join").String()
			db.Update(i, collection, data)
			return "success updated"
		}
	}
	return "nothing to update"
}

// TODO updateMany update document data
func updateMany(query string) (result string) {
	collection := gjson.Get(query, "collection").String() // + slash
	if collection == "" {
		return "ERROR! select no collection "
	}

	filter := gjson.Get(query, "filter").String()
	newData := gjson.Get(query, "data").String()

	tot := 0

	for i := 0; i <= db.lastId; i++ {
		if db.indexs[i].coll != collection {
			continue
		}

		data := db.Get(i, collection)

		if match(filter, data) {
			id := gjson.Get(data, "_id").Int()

			data = gjson.Get("["+data+","+newData+"]", "@join").String()
			db.Update(int(id), collection, data)
			tot++
		}
	}
	return str(tot) + " item updated"
}

// end
