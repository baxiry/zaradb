package engine

import "errors"

const (
	fOne   = "findOne"    // fo
	fMany  = "findMany"   // fm
	upOne  = "updateOne"  // uo
	upMany = "updateMany" // um
	dOne   = "deleteOne"  // do
	dMany  = "deleteMany" // dm

	//
	mtch    = "match"    // m
	sQery   = "subQuery" // sq
	orderby = "orderBy"  // ob
	fields  = "fields"   // f
	skip    = "skip"     // s
	limit   = "limit"    // l
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
