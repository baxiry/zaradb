package smap

type Map struct {
	key string
	val string
}

type List struct {
	list []Map
}

func NewSmap() *List {
	return &List{
		list: []Map{},
	}
}

func (l *List) Set(k, v string) {
	found := false
	for i := range l.list {
		if l.list[i].key == k {
			l.list[i].val = v
			found = true
			break
		}

	}
	if !found {
		l.list = append(l.list, Map{k, v})
	}
}

func (l *List) Get(k string) string {

	for i := 0; i < len(l.list); i++ {
		if l.list[i].key == k {
			return l.list[i].val
		}
	}
	return ""
}

func (l *List) Delete(k string) {
	list := []Map{}

	for i := 0; i < len(l.list); i++ {
		if l.list[i].key != k {
			list = append(list, l.list[i])
		}
	}
	l.list = list
}

func (l *List) Len() int {
	return len(l.list)
}
