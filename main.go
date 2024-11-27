package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

// Template struct implements the Echo Renderer interface
type Template struct {
	templates *template.Template
}

// Render method for the custom Template struct
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// Initialize SQLite database
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize db %v", err)
	}
	defer db.Close()

	runMigrations(db)
	seedDB(db)
	getReleases(db)

	// Initialize Echo
	e := echo.New()

	// Load templates
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	// Serve static files
	e.Static("/static", "static")

	// Define routes
	e.GET("/", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title": "Home Page",
		}
		return c.Render(http.StatusOK, "base.html", data)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8086"))
}
