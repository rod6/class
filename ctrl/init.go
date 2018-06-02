package ctrl

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "./class.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	// create table if not exists
	db.Exec(`create table if not exists courses(course varchar(100), teacher varchar(100));`)
	db.Exec(`create table if not exists classes(course varchar(100), student varchar(100));`)
	db.Exec(`create table if not exists absentees(course varchar(100), student varchar(100), memo varchar(100));`)
}

func Close() {
	DB.Close()
}
