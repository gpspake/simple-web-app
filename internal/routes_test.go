package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	e := echo.New()
	e.Renderer = &Template{TemplateDir: "./templates"}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	defer db.Close()

	createTestTables(db)
	seedTestReleases(db)
	seedTestArtists(db)
	seedTestReleaseArtists(db)
	populateReleasesFtsTable(db)
	SetupRoutes(e, db)

	t.Run("GET /", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Simple Web App")
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
		_, err := db.Exec("INSERT INTO releases (id, name, year) VALUES (?, ?, ?)",
			i, fmt.Sprintf("Album %d", i), startYear+(i-1))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed releases: %v", err))
		}
	}
}

func seedTestArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO artists (id, name) VALUES (?, ?)", i, fmt.Sprintf("Artist %d", i))
		if err != nil {
			panic(fmt.Sprintf("Failed to seed artists: %v", err))
		}
	}
}

func seedTestReleaseArtists(db *sql.DB) {
	for i := 1; i <= 30; i++ {
		_, err := db.Exec("INSERT INTO release_artists (id, release_id, artist_id) VALUES (?, ?, ?)", i, i, i)
		if err != nil {
			panic(fmt.Sprintf("Failed to seed release_artists: %v", err))
		}
	}
}

func populateReleasesFtsTable(db *sql.DB) {
	stmt, err := db.Prepare(`
		INSERT INTO releases_fts (release_id, artist_name, release_name, release_year)
		SELECT
			releases.id AS release_id,
			artists.name AS artist_name,
			releases.name AS release_name,
			releases.year AS release_year
		FROM
			release_artists
				JOIN
			artists ON release_artists.artist_id = artists.id
				JOIN
			releases ON release_artists.release_id = releases.id;
	`)
	if err != nil {
		panic(fmt.Sprintf("Failed to prepare statement: %v", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		panic(fmt.Sprintf("Failed to execute releases_fts query: %v", err))
	}
}
