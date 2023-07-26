package dblite

type Enginge struct{}

var pages = NewPages()

func (Enginge) Run() {
	pages.Open(RootPath)
}

func (Enginge) Stop() {
	pages.Close()
}

func NewEngine() *Enginge {
	return &Enginge{}
}
