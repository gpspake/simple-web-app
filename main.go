package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not connect to SQLite: %v", err)
	}

	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatalf("Could not create SQLite driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3",
		driver,
	)

	if err != nil {
		log.Fatalf("Could not initialize migrations: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not run migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}
