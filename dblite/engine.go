package dblite

type Enginge struct{}

var db = NewDatabase("test")

func (Enginge) Run() {

	// open all database pages
	db.Open()

	// check & init index map & firs page store
	initIndexs()
	println("ZARADB is runing on :" + PORT)
}

func (Enginge) Stop() {
	println("Enginge Closing...")
	db.Close()
}

func NewEngine() *Enginge {
	return &Enginge{}
}
