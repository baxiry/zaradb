package engine

import (
	"fmt"

	"github.com/tidwall/gjson"
)

var aggr Aggregate

// ok
func HandleQueries(query string) string {

	parsedQuery := gjson.Parse(query)

	switch parsedQuery.Get(action).Str {

	// aggregate actions
	case "aggregate":
		return aggr.aggrigate(parsedQuery)

	// io actions
	case "findOne":
		_, res := s.findOne(parsedQuery)
		return res

	case "findMany":
		return s.findMany(parsedQuery)

	case "findById":
		return s.findById(parsedQuery)

	case "insert":
		return s.insertOne(parsedQuery)

	case "insertMany":
		return s.insertMany(parsedQuery)

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
	case "createCollection":
		return s.createCollection(parsedQuery.Get(collection))

	case "deleteCollection":
		return s.deleteCollection(parsedQuery.Get(collection))

	case "getCollections":
		return s.getCollections()

	default:
		return fmt.Errorf("unknown '%s' cation", parsedQuery.Get("action").Str).Error()
	}
}
