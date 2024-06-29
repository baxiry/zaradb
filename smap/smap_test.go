package smap

import (
	"fmt"
	"testing"
)

type cse struct {
	key string
	val string
}

var cases = []cse{{"hi0", "hello0"}, {"hi1", "hello1"}, {"hi2", "hello2"}, {"hi3", "hello3"}, {"hi4", "hello4"}}

var smp = NewSmap()

func Test_Set(t *testing.T) {
	for i := 0; i < len(cases); i++ {
		smp.Set(cases[i].key, cases[i].val)
	}
}

func Test_Get(t *testing.T) {
	for i := 0; i < len(cases); i++ {
		v := smp.Get(cases[i].key)
		if v != cases[i].val {
			t.Errorf("have %s, want %s", v, cases[i].val)
		}
	}
}

func Test_Len(t *testing.T) {
	if smp.Len() != len(cases) {
		t.Errorf("have %d, want %d", smp.Len(), len(cases))
		//	t.Error(smp.list, cases)
	}
}

func Test_Set2(t *testing.T) {
	hello := "hello_hello"
	smp.Set("hi0", hello)
	v := smp.Get("hi0")
	if v != hello {
		t.Errorf("have %s, want %s", v, hello)
	}
}

func Test_Len2(t *testing.T) {
	if smp.Len() != len(cases) {
		//t.Errorf("have %d, want %d", smp.Len(), len(cases))
		t.Error("\n", smp.list, "\n", cases)
	}
}

func Test_Delete(t *testing.T) {
	hh := "hello_hello"
	smp.Delete("hi0")
	v := smp.Get("hi0")
	if v == hh {
		t.Errorf("have %s, want %s", v, hh)
	}
}

func Test_Len3(t *testing.T) {
	if smp.Len() >= len(cases) {
		t.Errorf("have %d, want %d", smp.Len(), len(cases)-1)
		//	t.Error("\n", smp.list, "\n", cases)
	}
}

// ---------------------- benchmark --------------------------------

var sm = NewSmap()
var m = map[string]string{}

func Benchmark_SMapSet(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			key := "key-" + fmt.Sprint(j)
			value := "value-" + fmt.Sprint(j)
			sm.Set(key, value)
		}
	}
}

func Benchmark_MapSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			key := "key-" + fmt.Sprint(j)
			value := "value-" + fmt.Sprint(j)
			m[key] = value
		}
	}
}

func Benchmark_SMapGet(b *testing.B) {

	var val string

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			key := "key-" + fmt.Sprint(j)
			val = sm.Get(key)
		}
	}
	_ = val
}

func Benchmark_MapGet(b *testing.B) {

	var val string
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			key := "key-" + fmt.Sprint(j)
			val = m[key]
		}
	}
	_ = val
}

/*
func main() {
	testing.Benchmark(BenchmarkSMapSet)
	testing.Benchmark(BenchmarkMapSet)
	testing.Benchmark(BenchmarkSMapGet)
	testing.Benchmark(BenchmarkMapGet)
}
*/
