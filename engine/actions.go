package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var aggr Aggregate

// ok
func HandleQueries(query string) string {

	parsedQuery := gjson.Parse(query)

	switch parsedQuery.Get("action").Str { // action

	// aggregate actions
	case "aggregate":
		return aggr.aggrigate(parsedQuery)

	case "count":
		// count(parsedQuery)
		return "count"

	case "sum":
		return "not implemented yet"

	case "avg":
		return "not implemented yet"

	case "min":
		return "not implemented yet"

	case "max":
		return "not implemented yet"

	// database actions
	case "findOne":
		return s.findOne(parsedQuery)

	case "findMany":
		return s.findMany(parsedQuery)

	case "findById":
		return s.findById(parsedQuery)

	case "insert":
		return s.insertOne(parsedQuery)

	case "insertMany":
		return s.insertMany(parsedQuery)

	// update
	case "updateById":
		return s.updateById(parsedQuery)

	case "updateOne":
		return s.updateOne(parsedQuery)

	case "updateMany":
		return s.updateMany(parsedQuery)

	case "deleteById":
		return s.deleteById(parsedQuery)

	case "deleteOne":
		return s.deleteOne(parsedQuery)

	case "deleteMany":
		return s.deleteMany(parsedQuery)

	case "transaction":
		return transaction(parsedQuery)

	// manage database
	case "create_collection":
		return s.createCollection(parsedQuery.Get("collection"))

	case "delete_collection":
		return s.deleteCollection(parsedQuery.Get("collection"))

	case "getCollections":
		//return showCollections(db.path)
		return s.getCollections()

	default:
		return fmt.Errorf("unknown '%s' cation", parsedQuery.Get("action").Str).Error()
	}
}
