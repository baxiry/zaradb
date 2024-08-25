package engine

import "errors"

const (
	// actions
	fOne   = "findOne"    // fo
	fMany  = "findMany"   // fm
	upOne  = "updateOne"  // uo
	upMany = "updateMany" // um
	dOne   = "deleteOne"  // do
	dMany  = "deleteMany" // dm

	//
	filter  = "match"    // m
	sQery   = "subQuery" // sq
	orderby = "orderBy"  // ob
	fields  = "fields"   // f
	skip    = "skip"     // s
	limit   = "limit"    // l

	// works with aggregation
	aggregate = "aggregate" // ag
	gmatch    = "gmatch"    //gm
	gskip     = "gskip"     // gs
	glimit    = "glimit"    // gl

	// separate strings
	siparator = "_:_"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
