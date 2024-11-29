package main

import (
	"database/sql"
	"fmt"
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

// Unique seeding functions for the test context
func testSeedReleases(db *sql.DB) {
	startYear := 1991
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO releases (id, name, year) VALUES (?, ?, ?)",
			i, fmt.Sprintf("Album %d", i), startYear+(i-1))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed releases: %v", err))
		}
	}
}

func testSeedArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO artists (id, name) VALUES (?, ?)", i, fmt.Sprintf("Artist %d", i))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed artists: %v", err))
		}
	}
}

func testSeedReleaseArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO release_artists (id, release_id, artist_id) VALUES (?, ?, ?)", i, i, i)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed release_artists: %v", err))
		}
	}
}
