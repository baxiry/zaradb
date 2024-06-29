package smap

func NewSmap() *list {
	return &list{
		list: []kv{},
	}
}

type kv struct {
	key string
	val interface{}
}

type list struct {
	list []kv
}

func (l *list) Set(k, v string) {
	found := false
	for i := range l.list {
		if l.list[i].key == k {
			l.list[i].val = v
			found = true
			break
		}
	}
	if !found {
		l.list = append(l.list, kv{k, v})
	}
}

func (l *list) Get(k string) string {

	for i := 0; i < len(l.list); i++ {
		if l.list[i].key == k {
			return l.list[i].val.(string)
		}
	}
	return ""
}
func (l *list) Delete(k string) {
	list := []kv{}

	for i := 0; i < len(l.list); i++ {
		if l.list[i].key != k {
			list = append(list, l.list[i])
		}
	}
	l.list = list
}

func (l *list) Len() int {
	return len(l.list)
}
