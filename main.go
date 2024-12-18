package main

import (
	"log"
	"os"
	"path/filepath"
	"simple-web-app/internal"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize SQLite database
	db, err := internal.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize db %v", err)
	}
	defer db.Close()

	internal.ResetDb()
	internal.RunMigrations(db)
	internal.SeedDB(db)

	// Initialize Echo
	e := echo.New()

	// Use environment variable or default to local development path
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		// Default path for local development
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		templateDir = filepath.Join(workingDir, "internal", "templates")
	}

	// Load templates
	e.Renderer = &internal.Template{TemplateDir: templateDir}

	internal.SetupRoutes(e, db)

	// Start server
	e.Logger.Fatal(e.Start(":8086"))
}
