package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
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

	connStr := "host=postgres_test user=testuser password=testpassword dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Apply migrations and seed data
	runMigrations(db)
	seedTestReleases(db)
	seedTestArtists(db)
	seedTestReleaseArtists(db)
	populateReleaseFts(db)
	SetupRoutes(e, db)

	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Home Page")
	})

	t.Run("GET /releases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/releases", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Album 1")
		assert.Contains(t, rec.Body.String(), "1991")
	})

	t.Run("Invalid Route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func seedTestReleases(db *sql.DB) {
	startYear := 1991
	for i := 1; i <= 30; i++ {
		_, err := db.Exec(
			"INSERT INTO releases (name, year) VALUES ($1, $2)",
			fmt.Sprintf("Album %d", i), startYear+(i-1),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed releases: %v", err))
		}
	}
}

func seedTestArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO artists (name) VALUES ($1)", fmt.Sprintf("Artist %d", i))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed artists: %v", err))
		}
	}
}

func seedTestReleaseArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO release_artists (release_id, artist_id) VALUES ($1, $2)", i, i)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed release_artists: %v", err))
		}
	}
}
