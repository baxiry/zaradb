package engine

import (
	"testing"
)

const input = `
[
   {
      "_id": 1,
      "name": "karam",
      "age": 16,
      "tal": 1.2
   },
   {
      "_id": 2,
      "name": "bxiry",
      "age": 16,
      "tal": 1.25
   },
   {
      "_id": 3,
      "name": "john",
      "age": 17,
      "tal": 1.65
   },
   {
      "_id": 4,
      "name": "adam",
      "age": 34,
      "tal": 1.63
   },
   {
      "_id": 5,
      "name": "hammam",
      "age": 33,
      "tal": 1.65
   },
   {
      "_id": 6,
      "name": "hisham",
      "age": 43,
      "tal": 1.66
   }
]`

const sortedStr = `[
   {
      "_id": 4,
      "name": "adam",
      "age": 34,
      "tal": 1.63
   },
   {
      "_id": 5,
      "name": "hammam",
      "age": 33,
      "tal": 1.65
   },
   {
      "_id": 6,
      "name": "hisham",
      "age": 43,
      "tal": 1.66
   },
   {
      "_id": 3,
      "name": "john",
      "age": 17,
      "tal": 1.65
   },
   {
      "_id": 1,
      "name": "karam",
      "age": 16,
      "tal": 1.2
   },
   {
      "_id": 1,
      "name": "karam",
      "age": 16,
      "tal": 1.2
   }
]`

func TestSort(t *testing.T) {
	// Test cases with different data and sort key
}
