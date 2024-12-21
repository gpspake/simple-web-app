package internal

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func InitializeApp() (*sql.DB, *echo.Echo) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	runMigrations(db)

	e := echo.New()

	// Use TEMPLATE_DIR environment variable if set, fallback to default path
	templateDir := os.Getenv("TEMPLATE_DIR")
	if templateDir == "" {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		templateDir = filepath.Join(workingDir, "internal", "templates")
	}

	e.Renderer = &Template{TemplateDir: templateDir}
	SetupRoutes(e, db)

	return db, e
}
