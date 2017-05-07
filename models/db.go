package models

import (
	"database/sql"
	"log"

	// Import "postgres" driver
	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initialize psql db connection
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
