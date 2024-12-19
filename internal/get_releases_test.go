package internal

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func TestGetReleases(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Seed the database
	testSeedReleases(db)
	testSeedArtists(db)
	testSeedReleaseArtists(db)
	populateReleaseFts(db)

	// Create an Echo logger
	e := echo.New()
	logger := e.Logger

	t.Run("Valid Limit and Offset", func(t *testing.T) {
		releases, err := getReleases(db, 5, 0, "", logger)
		if err != nil {
			t.Fatalf("Failed to fetch releases: %v", err)
		}

		if len(releases) != 5 {
			t.Errorf("Expected 5 releases, but got %d", len(releases))
		}
	})

	t.Run("Offset Exceeds Data", func(t *testing.T) {
		releases, err := getReleases(db, 5, 100, "", logger)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(releases) != 0 {
			t.Errorf("Expected 0 releases, but got %d", len(releases))
		}
	})

	t.Run("Invalid Limit", func(t *testing.T) {
		_, err := getReleases(db, -1, 0, "", logger)
		if err == nil {
			t.Fatalf("Expected error for invalid limit, but got nil")
		}
	})
}

func testSeedReleases(db *sql.DB) {
	startYear := 1991
	for i := 1; i <= 30; i++ {
		_, err := db.Exec(
			"INSERT INTO releases (name, year) VALUES ($1, $2)",
			fmt.Sprintf("Album %d", i), startYear+(i-1),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed releases: %v", err))
		}
	}
}

func testSeedArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO artists (name) VALUES ($1)", fmt.Sprintf("Artist %d", i))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed artists: %v", err))
		}
	}
}

func testSeedReleaseArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO release_artists (release_id, artist_id) VALUES ($1, $2)", i, i)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed release_artists: %v", err))
		}
	}
}
