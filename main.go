package main

import (
	"encoding/json"
	"fmt"
)

const version = "0.0.1"

type Person struct {
	Name    string
	Age     json.Number
	Job     string
	Contact string
	Address Address
}

type Adress struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println(err)
	}

	users := []Person{
		{"adam", "23", "09 534554432", "teacher", Address{"moroco", "sla", "harnose", "605555"}},
		{"jawad", "24", "09 534554432", "student", Address{"moroco", "cnatra", "jarnose", "305555"}},
		{"joha", "30", "09 534554432", "shef", Address{"moroco", "safro", "koarnose", "545555"}},
		{"jamous", "23", "09 534554432", "doctor", Address{"moroco", "safi", "saniya", "905555"}},
	}

	for value = range users {
		db.Write("users", value.Name, Person{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Job:     Value.Job,
			Adress:  value.Address,
		})
	}
	record, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(record)

	allusers := []Person{}

	for _, f := range records {
		userfound := Person{}
		if err := json.Unmarshal([]byte(f), &userfound); err != nil {
			fmt.Println(err)
		}
		allusers := append(allusers, userfound)
	}
	fmt.Println("all users is : ", allusers)

	if err := db.Delete("user", "john"); err != nil {
		fmt.Println(err)
	}

}
