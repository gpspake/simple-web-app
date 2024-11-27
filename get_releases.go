package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func getReleases(db *sql.DB) {
	// Query to select all rows from the 'releases' table
	rows, err := db.Query("SELECT id, name, year FROM releases")
	if err != nil {
		log.Fatalf("Failed to query releases: %v", err)
	}
	defer rows.Close()

	// Iterate over the rows and print the data
	fmt.Println("Releases:")
	for rows.Next() {
		var id int
		var name string
		var year int

		// Scan the columns into variables
		err := rows.Scan(&id, &name, &year)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}

		// Print the row
		fmt.Printf("ID: %d, Name: %s, Year: %d\n", id, name, year)
	}

	// Check for errors after iterating over rows
	if err = rows.Err(); err != nil {
		log.Fatalf("Error while iterating rows: %v", err)
	}
}
