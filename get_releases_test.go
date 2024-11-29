package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetReleases(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	createTestTables(db)
	seedTestReleases(db)
	seedTestArtists(db)
	seedTestReleaseArtists(db)

	t.Run("Valid Limit and Offset", func(t *testing.T) {
		releases, err := getReleases(db, 5, 0, nil)
		if err != nil {
			t.Fatalf("Failed to fetch releases: %v", err)
		}

		if len(releases) != 5 {
			t.Errorf("Expected 5 releases, but got %d", len(releases))
		}
	})

	t.Run("Offset Exceeds Data", func(t *testing.T) {
		releases, err := getReleases(db, 5, 100, nil)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(releases) != 0 {
			t.Errorf("Expected 0 releases, but got %d", len(releases))
		}
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		_, err := getReleases(db, -1, 0, nil)
		if err == nil {
			t.Fatalf("Expected error for invalid limit, but got nil")
		}
	})
}
