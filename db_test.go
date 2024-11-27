package main

import (
	"database/sql"
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
