package main

import (
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Initialize SQLite database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize db %v", err)
	}
	defer db.Close()

	resetDb()
	runMigrations(db)
	seedDB(db)

	// Initialize Echo
	e := echo.New()

	// Load templates
	e.Renderer = &Template{}

	setupRoutes(e, db)

	// Start server
	e.Logger.Fatal(e.Start(":8086"))
}
