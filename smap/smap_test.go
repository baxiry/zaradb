package smap

import (
	"testing"
)

type cse struct {
	key string
	val string
}

var cases = []cse{{"hi0", "hello0"}, {"hi1", "hello1"}, {"hi2", "hello2"}, {"hi3", "hello3"}, {"hi4", "hello4"}}

var smap = NewSmap()

func Test_Set(t *testing.T) {
	for i := 0; i < len(cases); i++ {
		smap.Set(cases[i].key, cases[i].val)
	}
}

func Test_Get(t *testing.T) {
	for i := 0; i < len(cases); i++ {
		v := smap.Get(cases[i].key)
		if v != cases[i].val {
			t.Errorf("have %s, want %s", v, cases[i].val)
		}
	}
}

func Test_Len(t *testing.T) {
	if smap.Len() != len(cases) {
		t.Errorf("have %d, want %d", smap.Len(), len(cases))
		//	t.Error(smap.list, cases)
	}
}

func Test_Set2(t *testing.T) {
	hello := "hello_hello"
	smap.Set("hi0", hello)
	if smap.Get("hi0") != hello {
		t.Errorf("have %s, want %s", smap.Get("hi"), hello)
	}
}

func Test_Len2(t *testing.T) {
	if smap.Len() != len(cases) {
		//t.Errorf("have %d, want %d", smap.Len(), len(cases))
		t.Error("\n", smap.list, "\n", cases)
	}
}

func Test_Delete(t *testing.T) {
	hh := "hello_hello"
	smap.Delete("hi0")
	if smap.Get("hi0") == hh {
		t.Errorf("have %s, want %s", smap.Get("hi0"), hh)
	}
}

func Test_Len3(t *testing.T) {
	if smap.Len() >= len(cases) {
		t.Errorf("have %d, want %d", smap.Len(), len(cases)-1)
		//	t.Error("\n", smap.list, "\n", cases)
	}
}

func BenchmarkMap(t *testing.B) {

	// keys := NewSmap()
}
