package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Initialize SQLite database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize db %v", err)
	}
	defer db.Close()

	runMigrations(db)
}
