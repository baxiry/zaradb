package smap

type List struct {
}

func NewSmap() *List {
	return &List{}
}

func (l *List) Set(k, v string) {

}

func (l *List) Get(k string) string {

	return ""
}

func (l *List) Delete(k string) {
}

func (l *List) Len() int {
	return 0
}
