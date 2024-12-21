package internal

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Verify the connection
	if err := db.Ping(); err != nil {
		t.Errorf("Failed to connect to PostgreSQL: %v", err)
	}
}

func TestResetDb(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	runMigrations(db)
	seedTestReleases(db)
	resetDb(db)

	// Verify tables are empty
	for _, table := range []string{"release", "artist", "release_artist"} {
		var count int
		err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to query table %s: %v", table, err)
		}
		if count != 0 {
			t.Errorf("Expected %s to be empty, but found %d rows", table, count)
		}
	}
}

func TestRunMigrations(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	runMigrations(db)

	// Verify migrations
	tables := []string{"release", "artist", "release_artist", "release_fts"}
	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = $1
			)`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			t.Fatalf("Failed to check table existence for %s: %v", table, err)
		}
		if !exists {
			t.Errorf("Expected table %s to exist, but it does not", table)
		}
	}
}

func TestSeedDB(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Run migrations and seed the database
	runMigrations(db)
	SeedDB(db)

	// Verify data in `release` table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM release").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query release table: %v", err)
	}
	if count != 30 {
		t.Errorf("Expected 30 releases, but got %d", count)
	}

	// Verify data in `artist` table
	err = db.QueryRow("SELECT COUNT(*) FROM artist").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query artist table: %v", err)
	}
	if count != 30 {
		t.Errorf("Expected 30 artists, but got %d", count)
	}

	// Verify data in `release_artist` table
	err = db.QueryRow("SELECT COUNT(*) FROM release_artist").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query release_artist table: %v", err)
	}
	if count != 30 {
		t.Errorf("Expected 30 release_artist, but got %d", count)
	}
}

func TestPopulateReleaseFts(t *testing.T) {
	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	runMigrations(db)
	seedTestReleases(db)
	seedTestArtists(db)
	seedTestReleaseArtists(db)
	PopulateReleaseFts(db)

	// Verify FTS table is populated
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM release_fts").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query release_fts: %v", err)
	}
	if count == 0 {
		t.Errorf("Expected release_fts to be populated, but found 0 rows")
	}
}
