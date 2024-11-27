package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestInitDB(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to connect to in-memory database: %v", err)
	}
}

func TestResetDb(t *testing.T) {
	// Define the file name used by resetDb
	const fileName = "data.db"

	// Step 1: Create a dummy file to simulate an existing database file
	t.Log("Creating dummy data.db file for test")
	file, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}
	file.Close()

	// Verify that the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Fatalf("Dummy data.db file was not created")
	}

	// Step 2: Call resetDb
	t.Log("Calling resetDb to delete and recreate the file")
	resetDb()

	// Verify that the file has been deleted and recreated
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("data.db file was not recreated")
	}

	// Clean up: Remove the file after the test
	t.Log("Cleaning up: Removing data.db")
	err = os.Remove(fileName)
	if err != nil {
		t.Fatalf("Failed to remove test data.db file: %v", err)
	}
}

func TestRunMigrations(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Run migrations
	runMigrations(db)

	// Verify if migrations applied correctly
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table'`)
	if err != nil {
		t.Fatalf("Failed to query tables after migrations: %v", err)
	}
	defer rows.Close()

	var tableCount int
	for rows.Next() {
		tableCount++
	}

	if tableCount == 0 {
		t.Errorf("No tables created after running migrations")
	}
}

func TestSeedDB(t *testing.T) {
	// Use an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create the 'releases' table (since seedDB assumes the table exists)
	createTableSQL := `
	CREATE TABLE releases (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		year INTEGER NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		t.Fatalf("Failed to create releases table: %v", err)
	}

	// Seed the database
	seedDB(db)

	// Query the releases table to verify if data has been seeded
	rows, err := db.Query("SELECT name, year FROM releases")
	if err != nil {
		t.Fatalf("Failed to query releases table: %v", err)
	}
	defer rows.Close()

	// Check the seeded data
	expectedReleasesCount := 30
	actualReleasesCount := 0
	expectedYear := 1991

	for rows.Next() {
		var name string
		var year int
		err := rows.Scan(&name, &year)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		// Validate the data
		if fmt.Sprintf("Album %d", actualReleasesCount+1) != name || year != expectedYear+actualReleasesCount {
			t.Errorf("Unexpected data: expected 'Album %d' for year %d, but got '%s' for year %d", actualReleasesCount+1, expectedYear+actualReleasesCount, name, year)
		}

		actualReleasesCount++
	}

	// Verify the correct number of records
	if actualReleasesCount != expectedReleasesCount {
		t.Errorf("Expected %d releases, but found %d", expectedReleasesCount, actualReleasesCount)
	}
}
