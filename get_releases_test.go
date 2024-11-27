package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetReleases(t *testing.T) {
	// Use an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create the 'releases' table
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

	// Seed the database with test data
	seedDB(db)

	// Redirect stdout to capture printed output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call getReleases
	getReleases(db)

	// Restore stdout and close the pipe writer
	w.Close()
	os.Stdout = oldStdout

	// Read the captured output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	r.Close()

	// Verify the output
	output := buf.String()
	expectedOutput := "Releases:\n"
	startYear := 1991

	for i := 1; i <= 30; i++ {
		expectedOutput += fmt.Sprintf("ID: %d, Name: Album %d, Year: %d\n", i, i, startYear+(i-1))
	}

	if output != expectedOutput {
		t.Errorf("Unexpected output:\nExpected:\n%s\nGot:\n%s", expectedOutput, output)
	}
}
