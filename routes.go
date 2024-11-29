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
		// Pass releases to the template
		data := map[string]interface{}{
			"Title": "Home Page",
		}

		// Render the template or return an error
		if err := c.Render(http.StatusOK, "index", data); err != nil {
			c.Logger().Errorf("Failed to render /: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	})

	e.GET("/about", func(c echo.Context) error {
		data := map[string]interface{}{
			"Title": "About",
		}
		return c.Render(http.StatusOK, "about", data)
	})

	e.GET("/releases", func(c echo.Context) error {
		// Read query parameters
		pageStr := c.QueryParam("page")
		limitStr := c.QueryParam("page_size")

		// Get releases with pagination
		releases, pagination, err := getPaginatedReleases(db, pageStr, limitStr, e.Logger, c.Request())

		if err != nil {
			e.Logger.Printf("Failed to get releases: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to load releases")
		}

		// Pass releases to the template
		data := map[string]interface{}{
			"Title":      "Releases",
			"Releases":   releases,
			"Page":       pageStr,
			"Pagination": pagination,
		}
		return c.Render(http.StatusOK, "releases", data)
	})
}
