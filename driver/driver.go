package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string
}

// NewDB create new directory to used as database
func NewDB(dir string) (*Driver, error) {
	dir = filepath.Clean(dir)

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
	}

	if _, err := os.Stat(dir); err != nil {
		log.Printf("using %s (database) already exists \n", dir)
		return &driver, nil
	}
	log.Printf("creating new databse at %s", dir)
	return &driver, os.MkdirAll(dir, 0744)
}

// Write write data to file.json
func (d *Driver) Write(collection, resource string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection - no place to save record, ")
	}
	if resource == "" {
		return fmt.Errorf("missing resource - unable to save record (no name)! ")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	err := os.MkdirAll(dir, 0744)
	if err != nil {
		fmt.Println("Error is : ", err)
		return err
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	err = ioutil.WriteFile(tmpPath, b, 0666)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

// Read method read data from file.json
func (d *Driver) Read(collection, resource string, v interface{}) error {
	if collection == "" {
		return errors.New("missing collection - no place to save record! ")
	}

	if resource == "" {
		return errors.New("missing resource = unable to save record (no name)! ")
	}

	record := filepath.Join(d.dir, collection, resource)
	if _, err := state(record); err != nil {
		return err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)
}

// ReadAll
func (d *Driver) ReadAll(collection string) ([]string, error) {

	if collection == "" {
		return nil, errors.New("missing collection - unable to read! ")
	}

	dir := filepath.Join(d.dir, collection)

	if _, err := state(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)

	var records []string

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))
	}

	return records, nil
}

// Delete
func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := state(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to fine file or dir named %s  %s", path, err)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")

	}
	return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}

func state(path string) (fi os.FileInfo, err error) {

	fi, err = os.Stat(path)
	if os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
		if err != nil {
			log.Fatal(err)
		}
		return fi, nil
	}

	return fi, nil
}
