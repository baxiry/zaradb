package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// test db is exist ? db is just a directory

func TestNewDB(t *testing.T) {
	db, err := NewDB("./")
	if err != nil {
		fmt.Println(err)
	}

	_, err = os.Stat("./")
	if os.IsNotExist(err) {
		t.Errorf("db is not exist %v, %v", db, err)
	}

}

// we need some data to test all driver functions
type Person struct {
	Name    string
	Age     json.Number
	Job     string
	Contact string
	Address Address
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

// cases:
var users = []Person{
	{"adam", "23", "09 534554432", "teacher", Address{"moroco", "sla", "harnose", "605555"}},
	{"jawad", "24", "09 534554432", "student", Address{"moroco", "cnatra", "jarnose", "305555"}},
	{"joha", "30", "09 534554432", "shef", Address{"moroco", "safro", "koarnose", "545555"}},
	{"jamous", "23", "09 534554432", "doctor", Address{"moroco", "safi", "saniya", "905555"}},
	{"koko", "23", "09 534554432", "doctor", Address{"moroco", "safi", "saniya", "905555"}},
	{"dodo", "23", "09 534554432", "doctor", Address{"moroco", "safi", "saniya", "905555"}},
}

// Test Wrilte to db
func TestWrite(t *testing.T) {

	db, err := NewDB("./")
	if err != nil {
		fmt.Println(err)
	}

	_, err = os.Stat("users")
	if os.IsNotExist(err) {
		if err != nil {
			fmt.Println(err)
		}

		for _, value := range users {
			db.Write("users", value.Name, Person{
				Name:    value.Name,
				Age:     value.Age,
				Contact: value.Contact,
				Job:     value.Job,
				Address: value.Address,
			})
		}
	}

	files, err := ioutil.ReadDir("./users")
	if err != nil {
		log.Fatal(err)
	}

	if len(files) != len(users) {
		t.Errorf("something wrong")
	}
	fmt.Println("all ok")

	//_ = os.Remove("users/adam.json")
}

// TEST ReadAll func
func TestReadAll(*testing.T) {

	db, _ := NewDB("./")

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(records)

	allusers := []Person{}

	for _, f := range records {
		userfound := Person{}
		if err := json.Unmarshal([]byte(f), &userfound); err != nil {
			fmt.Println(err)
		}
		allusers = append(allusers, userfound)
	}
	//fmt.Println("all users is : ", allusers)

	//if err := db.Delete("users", "adam"); err != nil {
	//fmt.Println(err)
	//}

	newRecords, err := db.ReadAll("users")
	for _, f := range newRecords {
		fmt.Println(f)
	}
}
