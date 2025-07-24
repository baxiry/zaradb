package engine

import "errors"

const (
	dbName     = "db"         // d
	collection = "collection" // c
	action     = "action"     // c

	// actions
	findOne    = "findOne"    // fo
	findMany   = "findMany"   // fm
	upOne      = "updateOne"  // uo
	upMany     = "updateMany" // um
	deleteOne  = "deleteOne"  // do
	deleteMany = "deleteMany" // dm
	aggregate  = "aggregate"  // ag

	//
	filter  = "match"    // m
	subQery = "subQuery" // sq
	fields  = "fields"   // f
	sort    = "sort"     // ob
	skip    = "skip"     // s
	limit   = "limit"    // l

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
