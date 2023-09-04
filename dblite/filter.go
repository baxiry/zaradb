package dblite

import "github.com/tidwall/gjson"

// Eq check numbers are equal
func Eq(json, field string, input int64) (result bool) {
	return input == gjson.Get(json, field).Int()
}

/*
// sNe check strings args are not equal
func sNe(json, field, match string) (result bool) {
	return match != gjson.Get(json, field).String()
}

// sEq check strings args is equal ?
func sEq(json, field, match string) (result bool) {
	return match == gjson.Get(json, field).String()
}



// Eq check numbers are not equal
func Ne(json, field, match string) (result bool) {
	return match != gjson.Get(json, field).String()
}

// Eq check if field is Greater then input arg
func Gt(json, field, match string) (result bool) {
	return match == gjson.Get(json, field).String()
}

// Eq check if field is Less then input arg
func Lt(json, field, match string) (result bool) {
	return match != gjson.Get(json, field).String()
}

// Eq check if field is Greater or Equal to input arg
func Ge(json, field, match string) (result bool) {
	return match == gjson.Get(json, field).String()
}

// Eq check if field is Less or Equal to input arg
func Le(json, field, match string) (result bool) {
	return match != gjson.Get(json, field).String()
}

*/
