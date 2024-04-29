package main

import "fmt"

// var query = `{"name":"adam","age": 21, "email": "jamal@email.com", "tele":"002345678"}}` //,"contact":{"home":12, "work":45, "tele":{"first":111, "Second":222}}}}`)
// var query = `{"name":"adam","age": {"$lte":30, "$gte":18}, "email": "jamal@email.com"}, "tele":"001234567"}}` //,"contact":{"home":12, "work":45, "tele":{"first":111, "Second":222}}}}`)
var query = `{"age": {"$lte":30, "$gte":18}}` //,"contact":{"home":12, "work":45, "tele":{"first":111, "Second":222}}}}`)

var data_ = `{"name":"adam","age": 18, "email": "jamal@email.com", "tele":"001234567"}}`

func test() {

	result := make([]string, 0)
	for _, v := range data {

		if match(query, v) {

			result = append(result, v)
		}
	}

	for k, v := range result {
		fmt.Println(k, v)
	}
	fmt.Println("len result is ", len(result))
	//res := valid(query, data_)
	//fmt.Println("\nres : ", res)

}

var data = []string{
	`{"name":"adam","age": 18, "email": "jamal@email.com", "tele":"001234567"}}`,
	`{"name":"karim", "age": 37, "contact":{"email": "karim@email.com", "tele":002345678 }}`,
	`{"name":"jawad", "age": 24, "contact":{"email": "jawad@email.com", "tele":001234567 }}`,
	`{"name":"imane", "age": 26, "contact":{"email": "imane@email.com", "tele":002345030 }}`,
	`{"name":"adams", "age": 31, "contact":{"email": "adams@email.com", "tele":002345000 }}`,
	`{"name":"johns", "age": 11, "contact":{"email": "adams@email.com", "tele":002345010 }}`,
	`{"name":"jawar", "age": 10, "contact":{"email": "jawar@email.com", "tele":002345022 }}`,
	`{"name":"kadir", "age": 13, "contact":{"email": "kadir@email.com", "tele":002345033 }}`,
	`{"name":"hamid", "age": 16, "contact":{"email": "hamid@email.com", "tele":002345004 }}`,
	`{"name":"rajab", "age": 18, "contact":{"email": "rajab@email.com", "tele":002345005 }}`,
	`{"name":"samir", "age": 19, "contact":{"email": "samir@email.com", "tele":002345006 }}`,
	`{"name":"akram", "age": 21, "contact":{"email": "akram@email.com", "tele":002345007 }}`,
	`{"name":"rabih", "age": 31, "contact":{"email": "rabih@email.com", "tele":002345008 }}`,
	`{"name":"monir", "age": 31, "contact":{"email": "monir@email.com", "tele":002345009 }}`,
}
