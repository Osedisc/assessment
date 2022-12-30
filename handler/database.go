package handler

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase(URL string) {
	var err error
	DB, err = sql.Open("postgres", URL)
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `
	CREATE TABLE IF NOT EXISTS expenses ( 
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err = DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
}
