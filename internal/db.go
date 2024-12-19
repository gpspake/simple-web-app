package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}

func resetDb(db *sql.DB) {
	tables := []string{"releases_fts", "release_artists", "artists", "releases"}

	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_name = $1
			)`
		err := db.QueryRow(query, table).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check table existence for %s: %v", table, err)
		}

		if exists {
			_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table))
			if err != nil {
				log.Fatalf("Failed to truncate table %s: %v", table, err)
			}
			log.Printf("Truncated table: %s", table)
		} else {
			log.Printf("Table %s does not exist. Skipping.", table)
		}
	}

	log.Println("Database reset successfully!")
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not create PostgreSQL driver: %v", err)
	}

	// Use environment variable to determine the migrations path
	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		log.Fatal("MIGRATIONS_PATH environment variable is not set")
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Could not initialize migrations: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not run migrations: %v", err)
	}

	log.Println("Migrations applied successfully!")
}

func seedReleases(db *sql.DB) {
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

	stmt, err := db.Prepare("INSERT INTO releases (name, year) VALUES ($1, $2)")
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

func seedArtists(db *sql.DB) {
	// Seed data for the 'artists' table
	type Artist struct {
		Name string
	}

	var artists = []Artist{
		{Name: "Queen"},
		{Name: "Radio"},
		{Name: "Eagle"},
		{Name: "Blurb"},
		{Name: "Cream"},
		{Name: "Oasis"},
		{Name: "Panic"},
		{Name: "Drake"},
		{Name: "Kyuss"},
		{Name: "Spark"},
		{Name: "Patti"},
		{Name: "Siren"},
		{Name: "Beach"},
		{Name: "Ratat"},
		{Name: "Reign"},
		{Name: "Shins"},
		{Name: "Smoke"},
		{Name: "Tracy"},
		{Name: "Peach"},
		{Name: "Moody"},
		{Name: "Suede"},
		{Name: "Flume"},
		{Name: "Tonic"},
		{Name: "Lorde"},
		{Name: "Exile"},
		{Name: "Mecca"},
		{Name: "Jewel"},
		{Name: "Spoon"},
		{Name: "Adele"},
		{Name: "Janes"},
	}

	stmt, err := db.Prepare("INSERT INTO artists (name) VALUES ($1)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, artist := range artists {
		_, err := stmt.Exec(artist.Name)
		if err != nil {
			log.Printf("Failed to insert artist '%s': %v", artist.Name, err)
		} else {
			log.Printf("Successfully inserted artist: '%s'", artist.Name)
		}
	}

	log.Println("Seeded Artists")
}

func seedReleaseArtists(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	log.Println("Preparing statement for release_artists")
	stmt, err := tx.Prepare("INSERT INTO release_artists (release_id, artist_id) VALUES ($1, $2)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	log.Println("Prepared statement successfully")

	for i := 1; i <= 30; i++ {
		_, err := stmt.Exec(i, i)
		if err != nil {
			log.Printf("Failed to insert row %d: %v", i, err)
		} else {
			fmt.Printf("Inserted row: release_id=%d, artist_id=%d\n", i, i)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Seeded release_artists")
}

func populateReleaseFts(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	// Clear the table to avoid duplicates
	_, err = tx.Exec("TRUNCATE releases_fts")
	if err != nil {
		log.Fatalf("Failed to truncate releases_fts: %v", err)
	}

	query := `
		INSERT INTO releases_fts (release_id, release_name, release_year, artist_name, tsvector_column)
		SELECT
			releases.id AS release_id,
			releases.name AS release_name,
			releases.year AS release_year,
			artists.name AS artist_name,
			to_tsvector(releases.name || ' ' || artists.name || ' ' || releases.year::TEXT) AS tsvector_column
		FROM
			release_artists
		JOIN releases ON release_artists.release_id = releases.id
		JOIN artists ON release_artists.artist_id = artists.id;
	`

	_, err = tx.Exec(query)
	if err != nil {
		log.Fatalf("Failed to populate releases_fts: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Populated Release FTS")
}

func SeedDB(db *sql.DB) {
	seedReleases(db)
	seedArtists(db)
	seedReleaseArtists(db)
	populateReleaseFts(db)
}
