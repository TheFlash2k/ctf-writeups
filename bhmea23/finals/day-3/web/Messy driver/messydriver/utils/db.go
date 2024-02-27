package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(dbpath string) *sql.DB {
	var err error
	db, err := sql.Open("sqlite3", dbpath)

	if err != nil {
		panic(err)
	}

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[+] Database connected: v%s\n", version)

	return db
}
