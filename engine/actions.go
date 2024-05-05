package engine

import (
	"github.com/tidwall/gjson"
)

// ok
func HandleQueries(query string) string {
	switch gjson.Get(query, "action").String() {

	// database actions
	case "findOne":
		return db.findOne(query)

	case "findMany":
		return db.findMany(query)

	case "findById":
		return db.findById(query)

	case "insert":
		return db.insertOne(query)

	case "insertMany":
		return db.insertMany(query)

	// update
	case "updateById":
		return db.updateById(query)

	case "updateOne":
		return db.updateOne(query)

	case "updateMany":
		return db.updateMany(query)

	case "deleteById":
		return db.deleteById(query)

	case "deleteOne":
		return db.deleteOne(query)

	case "deleteMany":
		return db.deleteMany(query)

	case "transaction":
		return transaction(query)

	// manage database
	case "create_collection":
		return createCollection(query)

	case "delete_collection":
		return deleteCollection(query)

	case "show_collection":
		//return showCollections(db.path)
		return "ont emplements yet"
	default:
		return "unknowen action"
	}
}
