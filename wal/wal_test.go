package wal

import (
	"os"
	"testing"
)

func TestWAL_WriteAndRead(t *testing.T) {
	os.RemoveAll("testwal")
	w, err := NewWAL("testwal")
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// أضف عددًا كافيًا لإجبار التدوير
	for i := 0; i < 10000; i++ {
		err := w.Write(Entry{
			TableID: byte(i % 3),
			Index:   uint64(i),
			Data:    []byte("hello world " + string(rune('A'+(i%26)))),
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	w.Flush()

	entries, err := w.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 10000 {
		t.Fatalf("expected 10000 entries, got %d", len(entries))
	}
}
