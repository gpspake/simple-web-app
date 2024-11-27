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
		data := map[string]interface{}{
			"Title": "Home Page",
		}
		return c.Render(http.StatusOK, "base.html", data)
	})
}
