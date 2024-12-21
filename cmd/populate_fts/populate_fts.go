package main

import (
	"log"
	"simple-web-app/internal"
)

func main() {
	db, err := internal.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	log.Println("Populating Full-Text Search table...")
	internal.PopulateReleaseFts(db)
	log.Println("Full-Text Search table population complete.")
}
