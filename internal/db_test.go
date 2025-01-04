package internal

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
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

	// Verify the database can accept queries
	_, err = db.Exec(`PRAGMA database_list;`)
	if err != nil {
		t.Errorf("Database is not operational: %v", err)
	}
}

func TestResetDb(t *testing.T) {
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
	ResetDb()

	// Verify the file is recreated
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.Errorf("data.db file was not recreated")
	}

	// Verify the recreated database is empty
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		t.Fatalf("Failed to open recreated database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("Failed to query tables in recreated database: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		t.Errorf("Expected no tables in the recreated database, but found some")
	}

	// Clean up
	os.Remove(fileName)
}

func TestRunMigrations(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Dynamically locate the migrations directory relative to this file
	_, testFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get test file path")
	}
	migrationsDir := filepath.Join(filepath.Dir(testFilePath), "../migrations")

	RunMigrations(db, migrationsDir)

	// Verify migrations
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		t.Fatalf("Failed to query tables after migrations: %v", err)
	}
	defer rows.Close()

	expectedTables := map[string]bool{
		"releases":        false,
		"artists":         false,
		"release_artists": false,
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Fatalf("Failed to scan table name: %v", err)
		}
		if _, ok := expectedTables[tableName]; ok {
			expectedTables[tableName] = true
		}
	}

	for table, created := range expectedTables {
		if !created {
			t.Errorf("Expected table %s was not created", table)
		}
	}
}

func TestSeedDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	createTestTables(db)
	SeedDB(db)

	// Verify data in `releases` table
	rows, err := db.Query("SELECT COUNT(*) FROM releases")
	if err != nil {
		t.Fatalf("Failed to query releases table: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			t.Fatalf("Failed to scan releases count: %v", err)
		}
	}

	if count != 30 {
		t.Errorf("Expected 30 releases, but got %d", count)
	}

	// Add similar assertions for `artists` and `release_artists` if necessary
}
