package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	// Initialize Echo
	e := echo.New()

	// Register the renderer
	e.Renderer = &Template{}

	// Use an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer db.Close()

	// Set up the database schema for testing
	_, err = db.Exec(`CREATE TABLE releases (id INTEGER PRIMARY KEY, name TEXT, year INTEGER)`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	_, err = db.Exec(`INSERT INTO releases (name, year) VALUES 
		('Album 1', 1991), 
		('Album 2', 1992)`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Set up routes
	setupRoutes(e, db)

	// Test the home route
	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Welcome to the Home Page")
	})

	// Test the about route
	t.Run("GET /about", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "This is the about page.")
	})

	// Test the releases route
	t.Run("GET /releases", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/releases", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Album 1")
		assert.Contains(t, rec.Body.String(), "1991")
	})
}
