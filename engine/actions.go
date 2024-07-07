package engine

import (
	"github.com/tidwall/gjson"
)

// ok
func HandleQueries(query string) string {

	parsedQuery := gjson.Parse(query)

	switch parsedQuery.Get("action").String() { // action
	// database actions
	case "findOne":
		return db.findOne(parsedQuery)

	case "findMany":
		return db.findMany(parsedQuery)

	case "findById":
		return db.findById(parsedQuery)

	case "insert":
		return db.insertOne(parsedQuery)

	case "insertMany":
		return db.insertMany(parsedQuery)

	// update
	case "updateById":
		return db.updateById(parsedQuery)

	case "updateOne":
		return db.updateOne(parsedQuery)

	case "updateMany":
		return db.updateMany(parsedQuery)

	case "deleteById":
		return db.deleteById(parsedQuery)

	case "deleteOne":
		return db.deleteOne(parsedQuery)

	case "deleteMany":
		return db.deleteMany(parsedQuery)

	case "transaction":
		return transaction(parsedQuery)

	// manage database
	case "create_collection":
		return createCollection(parsedQuery.Get("collection"))

	case "delete_collection":
		return deleteCollection(parsedQuery.Get("collection"))

	case "getCollections":
		//return showCollections(db.path)
		return getCollections()

	// trying sqlite query
	case "sql":
		return querySql(parsedQuery)

	default:
		return "unknowen action"
	}
}

// hmmmmm sql
func querySql(query gjson.Result) string {
	qr := query.Get("query").Str
	res, _ := db.db.Query(qr)
	record := ""
	result := "["
	for res.Next() {
		res.Scan(&record)
		result += record + ","
	}
	return result[:len(result)-1] + "]"
}
