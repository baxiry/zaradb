package dblite

import (
	"fmt"
	"io"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// data enginge

// At is where enginge insert data in page
var At int
var MaxObjects int64 = 1_000

// nuber page to make namePage from id
var numberPage int64 = 0

const slash = string(os.PathSeparator) // not tested for windos

// var namePage := 0
// Insert
func Insert(query string) (res string) {

	collection := gjson.Get(query, "in").String() + slash
	if len(collection) == 1 {
		return fmt.Sprint("failure insert. insert into no collection")
	}

	pName := db.PrimaryIndex / MaxObjects // page name as int
	// TODO check here . my be a bug
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
	AppendIndex(db.Pages[db.Name+collection+pi], At, size)

	//UpdateIndex(db.Pages[db.Name+collection+pi], int(id), int64(At), int64(size))

	At += size
	db.PrimaryIndex++
	return fmt.Sprint("Success Insert, _id: ", db.PrimaryIndex-1)
}

// Update update document data
func Update(query string) (result string) {
	collection := gjson.Get(query, "in").String() + slash
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
	size := len(data)

	UpdateIndex(db.Pages[db.Name+collection+pi], int(id), int64(At), int64(size))

	//		AppendIndex(db.Pages[db.Name+collection+pi], At, size)

	At += size

	return "Success update"
}

// delete
func DeleteById(query string) (result string) {

	id := gjson.Get(query, "_id").Int()
	in := gjson.Get(query, "in").String() + slash

	path := db.Name + in + fmt.Sprint(db.PrimaryIndex/MaxObjects)

	fmt.Println("path id DeleteById: ", path)

	UpdateIndex(db.Pages[path], int(id), 0, 0)

	//fmt.Println(IndexsCache.indexs)
	IndexsCache.indexs[id] = [2]int64{0, 0}
	//fmt.Println(IndexsCache.indexs)

	return "Delete Success!"
}

// Select reads data form docs
func SelectById(query string) (result string) {
	id := gjson.Get(query, "where_id").Int()
	if int(id) >= len(IndexsCache.indexs) {

		iLog.Println("no found index", id)
		return fmt.Sprintf("Not Found _id %v\n", id)
	}

	at := IndexsCache.indexs[id][0]
	size := IndexsCache.indexs[id][1]

	in := gjson.Get(query, "in").String() + slash
	//fmt.Println("table is : ", in)
	// TODO check from if exist!

	path := db.Name + in + fmt.Sprintf("%d", id/MaxObjects)

	result = Get(db.Pages[path], at, int(size))

	return result
}

// appends data to Pagefile & returns file size or error
func Append(file *os.File, data string) (size int, err error) {
	size, err = file.WriteAt([]byte(data), int64(At))
	if err != nil {
		eLog.Println("Error WriteString ", err)
	}
	return size, err
}

// Select reads data form docs
func Select(filter string) (result string) {
	id := gjson.Get(filter, "_id").String()
	fmt.Println("id is ", id)

	return result
}

// gets data from *file, takes at (location) & buffer size
func Get(file *os.File, at int64, size int) string {

	buffer := make([]byte, size)

	// read at
	n, err := file.ReadAt(buffer, at)
	if err != nil && err != io.EOF {
		eLog.Println(err)
		eLog.Println("At", at)
		eLog.Println("Size", size)
		return "ERROR form Get::ReadAt func"
	}

	// out the buffer content
	return string(buffer[:n])
}

// Delete removes document
func Delete(path string) (err error) {
	return
}

func selectFields(query string) string {
	return ""
}

// wht is fast ? remander or divider ? 3000/10 or 3000 % 10. for speed during extract dataPage form id

/*

// Update update document data
func Update(query string) (result string) {
	collection := gjson.Get(query, "in").String() + slash
	if len(collection) == 1 {
		return "ERROR! select no collection "
	}

	data := SelectById(query)
	newData := gjson.Get(query, "data").String()

	// `{"object":{"first":1,"second":2,"third":3}}`
	jsonParsed, err := gabs.ParseJSON([]byte(newData))
	if err != nil {
		return fmt.Sprintf("ERROR: parse data json %s", err)
	}

	// extract fields that need to update
	for field, val := range jsonParsed.ChildrenMap() {
		result, _ = sjson.Set(data, field, val)
		data = result
	}

	id := gjson.Get(data, "_id").Int()

	path := db.Name + collection + fmt.Sprint(id/MaxObjects)

	_, err = Append(db.Pages[path], data)
	if err != nil {
		return fmt.Errorf("ERROR! from Append %v\n", err).Error()
	}

	// Update index
	size := len(data)

	UpdateIndex(db.Pages[db.Name+collection+pi], int(id), int64(At), int64(size))

	At += size

	return "Success update"
}

*/
