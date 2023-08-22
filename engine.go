package dblite

type Enginge struct{}

var db = NewDatabase("test")

func (Enginge) Run() {

	// println("Engine Runing...")
	db.Open()

	// check & init index map & firs page store
	initIndexsFile()

	//db := NewDatabase()
	initIndex()

}

func (Enginge) Stop() {
	println("Enginge Closing...")
	db.Close()
}

func NewEngine() *Enginge {
	return &Enginge{}
}
