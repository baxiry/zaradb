package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	db = initDatabase()
	defer db.Close()

	loopping(db)

	http.HandleFunc("/info", info)
	http.HandleFunc("/time", expiration)
	http.HandleFunc("/new", newBoot)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", deleteBoot)
	http.HandleFunc("/black", blackMed)

	//http.HandleFunc("/auth", auth)
	http.ListenAndServe(":8001", nil)
}

func blackMed(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println("one=" + r.FormValue("one"))
	log.Println("two=" + r.FormValue("two"))
	fmt.Fprintf(w, "Gorilla!\n")
}

// time expiration license
func expiration(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	serial := url.Get("serial")

	ts := ""
	row := db.QueryRow("select unix_timestamp(ts) from licenses.boots where serial=?", serial)

	if err := row.Scan(&ts); err != nil {
		fmt.Println("ts is : ", ts)
		ErrorCheck(err)
	}

	times, err := strconv.Atoi(ts)
	ErrorCheck(err)
	nw := int(time.Now().Unix())
	left := 31 - (nw-times)/60/60/24

	fmt.Fprintf(w, "%d", left)
}

// TODO check this auth func
// check by serial
func info(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	serial := url.Get("serial")

	row := db.QueryRow("select name, ipaddress, ts from licenses.boots where serial=?", serial)

	var name, ipaddress, time string
	if err := row.Scan(&name, &ipaddress, &time); err != nil {
		ErrorCheck(err)
	}

	fmt.Fprintf(w, "%s, %s, %s", name, ipaddress, time)
}

// deleteBoot
func deleteBoot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	serial := url.Get("serial")

	// delete boot from db by serial
	stmt, e := db.Prepare("delete from licenses.boots where serial=?")
	ErrorCheck(e)

	_, e = stmt.Exec(serial)
	ErrorCheck(e)

	fmt.Fprintf(w, "serial : %s\n", serial)
}

func newBoot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	serial := url.Get("serial")
	name := url.Get("name")
	ipaddress := url.Get("ip")

	// prepare
	stmt, e := db.Prepare("insert into licenses.boots(name, serial, ipaddress ) values (?, ?,?)")
	ErrorCheck(e)

	//execute
	_, err := stmt.Exec(name, serial, ipaddress) //,serial())
	ErrorCheck(e)
	if err != nil {
		fmt.Fprintf(w, "wrong")
		return
	}

	fmt.Fprintf(w, "بوت جديد\nname: %s\nserial : %s\nipaddress %s", name, serial, ipaddress)
}

// changeIpAddr update ip addr
func update(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	serial := url.Get("serial")
	ipaddress := url.Get("ip")
	name := url.Get("name")

	fmt.Printf("name: %s\nip: %s\nserial: %s\n", name, ipaddress, serial)

	// update name or ip or both
	var column, value string
	if len(name) > 1 && len(ipaddress) > 5 {
		stmt, e := db.Prepare("update licenses.boots set name=?, ipaddress=? where serial=?")
		ErrorCheck(e)
		_, e = stmt.Exec(name, ipaddress, serial)
		ErrorCheck(e)
		fmt.Fprintf(w, "update serial : %s\nipaddress %s", serial, ipaddress)
		return
	}

	if len(name) > 1 && len(ipaddress) < 3 {
		column = "name"
		value = name
	} else if len(name) < 1 && len(ipaddress) > 3 {
		column = "ipaddress"
		value = ipaddress
	} else {
		fmt.Fprintf(w, "nothing to update\n you messing name ipaddress")
		return
	}

	stmt, e := db.Prepare("update licenses.boots set " + column + "=? where serial=?")
	ErrorCheck(e)
	_, e = stmt.Exec(value, serial)
	ErrorCheck(e)
	fmt.Fprintf(w, "update %s : %s for serial : %s\n", column, value, serial)
}

// initialaze database
func initDatabase() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8&parseTime=True&loc=Local")
	ErrorCheck(err)

	return db
}

func ErrorCheck(err error) {
	if err != nil {
		println(err.Error())
	}
}

// loop ping for active db connextion
func loopping(db *sql.DB) {
	go func() {
		for {
			err := db.Ping()
			ErrorCheck(err)
			time.Sleep(time.Minute * 1)
		}
	}()
}

/*
func auth(w http.ResponseWriter, r *http.Request) {
	// check boot serial if not run on aother ip addres
	url := r.URL.Query()
	serial := url.Get("serial")
	fmt.Println(serial)

	time := ""
	row := db.QueryRow("select ts from licenses.boots where serial=?", serial)

	if err := row.Scan(&time); err != nil {
		ErrorCheck(err)
	}

	fmt.Print(time)

	fmt.Fprintf(w, time)
}

*/

//fs := http.FileServer(http.Dir("static/"))
//http.Handle("/static/", http.StripPrefix("/static/", fs))
