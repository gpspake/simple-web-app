package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func initDB() (*sql.DB, error) {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

func resetDb() {
	// Define the file name
	const fileName = "data.db"

	// Check if the file exists
	if _, err := os.Stat(fileName); err == nil {
		// If it exists, delete it
		log.Println("Deleting existing data.db...")
		if err := os.Remove(fileName); err != nil {
			log.Fatalf("Error deleting file: %v", err)
			return
		}
		log.Println("File deleted successfully.")
	} else if !os.IsNotExist(err) {
		// Handle other errors from os.Stat
		log.Fatalf("Error checking file: %v", err)
		return
	}

	// Recreate the file
	log.Println("Recreating data.db...")
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	log.Println("File created successfully.")
}

func runMigrations(db *sql.DB) {
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

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not run migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}

func seedDB(db *sql.DB) {
	// Seed data for the 'releases' table
	var releases []struct {
		Name string
		Year int
	}

	startYear := 1991
	for i := 1; i <= 30; i++ {
		releases = append(releases, struct {
			Name string
			Year int
		}{
			Name: fmt.Sprintf("Album %d", i),
			Year: startYear + (i - 1),
		})
	}

	// Prepare the INSERT statement
	stmt, err := db.Prepare("INSERT INTO releases (name, year) VALUES (?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Insert data into the table
	for _, release := range releases {
		_, err := stmt.Exec(release.Name, release.Year)
		if err != nil {
			log.Printf("Failed to insert release '%s': %v", release.Name, err)
		} else {
			log.Printf("Successfully inserted release: '%s'", release.Name)
		}
	}

	log.Println("Seeded releases")
}
