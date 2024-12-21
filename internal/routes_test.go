package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
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
	PopulateReleaseFts(db)
	SetupRoutes(e, db)

	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Home Page")
	})

	t.Run("GET / - Rendering Error", func(t *testing.T) {
		// Simulate a rendering error by using a faulty renderer
		e.Renderer = &FaultyRenderer{}
		defer func() { e.Renderer = &Template{TemplateDir: templateDir} }() // Restore the original renderer

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "forced rendering error") // Adjusted expectation
	})

	t.Run("GET /about - Rendering Error", func(t *testing.T) {
		e.Renderer = &FaultyRenderer{}
		defer func() { e.Renderer = &Template{TemplateDir: templateDir} }()

		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Internal Server Error") // Adjusted expectation
	})

	t.Run("GET /releases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/releases", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Album 1")
		assert.Contains(t, rec.Body.String(), "1991")
	})

	t.Run("GET /releases - HTMX Partial", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/releases", nil)
		req.Header.Set("HX-Request", "true")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Album 1")
		assert.Contains(t, rec.Body.String(), "1991")
	})

	t.Run("GET /artist/:artist_id - Non-existent Artist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/artist/9999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "artist not found")
	})

	t.Run("GET /release/:release_id - Invalid ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/release/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid release ID")
	})

	t.Run("GET /release/:release_id - Non-Existent Release", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/release/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "release not found")
	})

	t.Run("GET /artist/:artist_id - Valid Artist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/artist/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Queen")
	})

	t.Run("GET /artist/:artist_id - Invalid Artist ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/artist/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid artist ID")
	})

	t.Run("GET /artist/:artist_id - Non-existent Artist", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/artist/999", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "artist not found")
	})

	t.Run("Invalid Route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/invalid", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// FaultyRenderer simulates a rendering error
type FaultyRenderer struct{}

// Render forces a rendering error
func (r *FaultyRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return errors.New("forced rendering error")
}

// BrokenDB simulates a database error
type BrokenDB struct{}

func (b *BrokenDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New("forced query error")
}

func (b *BrokenDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{} // Simulate empty row
}

func (b *BrokenDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New("forced database error")
}

func seedTestReleases(db *sql.DB) {
	startYear := 1991
	for i := 1; i <= 30; i++ {
		_, err := db.Exec(
			"INSERT INTO release (title, year) VALUES ($1, $2)",
			fmt.Sprintf("Album %d", i), startYear+(i-1),
		)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed releases: %v", err))
		}
	}
}

func seedTestArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO artist (name) VALUES ($1)", fmt.Sprintf("Artist %d", i))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed artists: %v", err))
		}
	}
}

func seedTestReleaseArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO release_artist (release_id, artist_id) VALUES ($1, $2)", i, i)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed release_artist: %v", err))
		}
	}
}
