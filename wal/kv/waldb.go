package wal

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const maxFileSize = 1024 * 1024 * 2 // 200 for 200MB

type Entry struct {
	TableID byte
	Index   uint64
	Data    []byte
}

type WAL struct {
	dir        string
	activeFile *os.File
	files      []*os.File
	currentID  int
	currentSz  int64
}

func NewWAL(dir string) (*WAL, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	w := &WAL{dir: dir}
	if err := w.createNewFile(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *WAL) createNewFile() error {
	w.currentID++
	filename := fmt.Sprintf("wal-%06d.log", w.currentID)
	fullPath := filepath.Join(w.dir, filename)
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if w.activeFile != nil {
		w.files = append(w.files, w.activeFile) // Keep previous open for reads
	}
	w.activeFile = f
	w.currentSz = 0
	return nil
}

func (w *WAL) Write(entry Entry) error {
	totalLen := 2 + 1 + 8 + len(entry.Data)
	if w.currentSz+int64(totalLen) >= maxFileSize {
		if err := w.createNewFile(); err != nil {
			return err
		}
	}

	buf := make([]byte, totalLen)
	binary.BigEndian.PutUint16(buf[0:2], uint16(len(entry.Data)))
	buf[2] = entry.TableID
	binary.BigEndian.PutUint64(buf[3:11], entry.Index)
	copy(buf[11:], entry.Data)

	n, err := w.activeFile.Write(buf)
	if err != nil {
		return err
	}
	w.currentSz += int64(n)
	return nil
}

func (w *WAL) Flush() error {
	return w.activeFile.Sync()
}

func (w *WAL) Close() error {
	for _, f := range append(w.files, w.activeFile) {
		f.Close()
	}
	return nil
}

func (w *WAL) ReadAll() ([]Entry, error) {
	allFiles := append(w.files, w.activeFile)
	var entries []Entry

	for _, f := range allFiles {
		_, err := f.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		for {
			header := make([]byte, 11)
			_, err := io.ReadFull(f, header)
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}

			dataLen := binary.BigEndian.Uint16(header[0:2])
			tableID := header[2]
			index := binary.BigEndian.Uint64(header[3:11])

			data := make([]byte, dataLen)
			_, err = io.ReadFull(f, data)
			if err != nil {
				return nil, err
			}

			entries = append(entries, Entry{
				TableID: tableID,
				Index:   index,
				Data:    data,
			})
		}
	}
	return entries, nil
}
