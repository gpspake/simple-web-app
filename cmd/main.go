package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"simple-web-app/internal"
)

func main() {
	// Initialize and configure the application
	db, e := internal.InitializeApp()

	defer db.Close()

	// Start the server
	e.Logger.Fatal(e.Start(":8086"))
}
