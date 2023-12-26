package dblite

import (
	"github.com/tidwall/gjson"
)

// ok
func HandleQueries(query string) string {
	switch gjson.Get(query, "action").String() {

	case "findOne":
		return findOne(query)

	case "findMany":
		return findMany(query)

	case "findById":
		return findById(query)

	case "insert":
		return insert(query)

	case "update":
		return update(query)

	case "deleteById":
		return deleteById(query)

	case "deleteOne":
		return deleteOne(query)

	case "deleteMany":
		return deleteMany(query)

	// manage database
	case "create_collection":
		return createCollection(query)

	case "delete_collection":
		return deleteCollection(query)

	case "show_collection":
		return showCollections(db.path)

	default:
		return "unknowen action"
	}
}
