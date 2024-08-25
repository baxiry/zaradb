package engine

import "errors"

const (
	dbname     = "db"         // d
	collection = "collection" // c

	// actions
	fOne      = "findOne"    // fo
	fMany     = "findMany"   // fm
	upOne     = "updateOne"  // uo
	upMany    = "updateMany" // um
	dOne      = "deleteOne"  // do
	dMany     = "deleteMany" // dm
	aggregate = "aggregate"  // ag

	//
	filter = "match"    // m
	sQery  = "subQuery" // sq
	sort   = "sort"     // ob
	fields = "fields"   // f
	skip   = "skip"     // s
	limit  = "limit"    // l

	// works with aggregation
	gmatch = "gmatch" // h
	gsort  = "gsort"  // gs
	gskip  = "gskip"  // gp
	glimit = "glimit" // gl

	// separate strings
	siparator = "-:-"
)

var (
	ErrDuplicate    = errors.New(" already exists ")
	ErrNotExists    = errors.New(" not exists ")
	ErrUpdateFailed = errors.New(" update failed ")
	ErrDeleteFailed = errors.New(" delete failed ")
	ErrIdNotExists  = errors.New("a group specification must include an _id")
)
