package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
)

func setupRoutes(e *echo.Echo, db *sql.DB) {
	// Serve static files
	e.Static("/static", "static")

	// Define routes
	e.GET("/", func(c echo.Context) error {
		releases, err := getReleases(db)

		if err != nil {
			e.Logger.Printf("Failed to get releases: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to load releases")
		}

		// Pass releases to the template
		data := map[string]interface{}{
			"Title":    "Home Page",
			"Releases": releases,
		}
		return c.Render(http.StatusOK, "base.html", data)
	})
}
