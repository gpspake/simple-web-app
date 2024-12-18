package internal

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
	"path/filepath"
)

func InitDB() (*sql.DB, error) {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}

func ResetDb() {
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

func RunMigrations(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatalf("Could not create SQLite driver: %v", err)
	}

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get current working directory: %v", err)
	}
	migrationsPath := filepath.Join(basePath, "migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
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

func seedReleases(db *sql.DB) {
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

	// Prepare the INSERT statement
	stmt, err := db.Prepare("INSERT INTO artists (name) VALUES (?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Insert data into the table
	for _, artist := range artists {
		_, err := stmt.Exec(artist.Name)
		if err != nil {
			log.Printf("Failed to insert artist '%s': %v", artist.Name, err)
		} else {
			log.Printf("Successfully inserted artist: '%s'", artist.Name)
		}
	}

	log.Println("Seeded Releases")
}

func seedReleaseArtists(db *sql.DB) {
	// Seed data for the 'release_artists' table
	type Artist struct {
		ReleaseId int
		ArtistId  int
	}

	// Seed the release_artists table
	tx, err := db.Begin() // Use a transaction for better performance
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare("INSERT INTO release_artists (id, release_id, artist_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for i := 1; i <= 30; i++ {
		_, err := stmt.Exec(i, i, i)
		if err != nil {
			log.Printf("Failed to insert row %d: %v", i, err)
		} else {
			fmt.Printf("Inserted row: id=%d, release_id=%d, artist_id=%d\n", i, i, i)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Seeded release_artists")
}

func populateReleaseFts(db *sql.DB) {

	// Populate 'release_fts' virtual table
	tx, err := db.Begin() // Use a transaction for better performance
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO releases_fts (release_id, artist_name, release_name, release_year)
		SELECT
			releases.id AS release_id,
			artists.name AS artist_name,
			releases.name AS release_name,
			releases.year AS release_year
		FROM
			release_artists
				JOIN
			artists ON release_artists.artist_id = artists.id
				JOIN
			releases ON release_artists.release_id = releases.id;
	`)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("Failed to execute releases_fts query %v", err)
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
