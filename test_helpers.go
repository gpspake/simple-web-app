package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func createTestTables(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	CREATE TABLE releases (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		year INTEGER NOT NULL
	);

	CREATE TABLE artists (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);

	CREATE TABLE release_artists (
		id INTEGER PRIMARY KEY,
		release_id INTEGER NOT NULL REFERENCES releases(id),
		artist_id INTEGER NOT NULL REFERENCES artists(id)
	);

	CREATE VIRTUAL TABLE releases_fts USING fts5
	(
		release_id UNINDEXED,
		release_name,
		release_year,
		artist_name,
		tokenize="trigram"
	);
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to execute schema: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Schema executed and transaction committed successfully")
	return nil
}

func cleanupTestDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		fmt.Printf("Failed to close test DB: %v\n", err)
	}
}
