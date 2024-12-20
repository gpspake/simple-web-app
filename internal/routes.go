package internal

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func SetupRoutes(e *echo.Echo, db DBQuerier) {
	// Serve static files
	e.Static("/static", "static")

	// Define routes
	e.GET("/", func(c echo.Context) error {
		// Pass releases to the template
		data := map[string]interface{}{
			"Title":        "Home Page",
			"CurrentRoute": c.Request().URL.Path,
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
			"Title":        "About",
			"CurrentRoute": c.Request().URL.Path,
		}
		return c.Render(http.StatusOK, "about", data)
	})

	e.GET("/releases", func(c echo.Context) error {
		// Read query parameters
		pageStr := c.QueryParam("page")
		limitStr := c.QueryParam("page_size")
		searchQuery := c.QueryParam("q")

		// Get releases with pagination and search
		releases, pagination, err := getPaginatedReleases(db, pageStr, limitStr, searchQuery, e.Logger, c.Request())

		if err != nil {
			e.Logger.Printf("Failed to get releases: %v", err)
			return c.String(http.StatusInternalServerError, "Failed to load releases")
		}

		// Render appropriate template (full page or HTMX partial)
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(http.StatusOK, "releases_partial", map[string]interface{}{
				"Releases":   releases,
				"Pagination": pagination,
			})
		}

		data := map[string]interface{}{
			"Title":        "Releases",
			"Releases":     releases,
			"Page":         pageStr,
			"Pagination":   pagination,
			"IncludeHTMX":  true,
			"CurrentRoute": c.Request().URL.Path,
		}

		return c.Render(http.StatusOK, "releases", data)
	})

	e.GET("/artist/:artist_id", func(c echo.Context) error {
		// Parse artist ID from the route
		artistID, err := strconv.Atoi(c.Param("artist_id"))
		if err != nil || artistID <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid artist ID"})
		}

		// Call getArtist to fetch artist details
		artist, err := getArtist(db, artistID, c.Logger())
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusNotFound, map[string]string{"error": "artist not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch artist details"})
		}

		// Prepare data for the template
		data := map[string]interface{}{
			"Title":        "Artist Details",
			"CurrentRoute": c.Request().URL.Path,
			"Artist":       artist, // Pass the artist struct to the template
		}

		return c.Render(http.StatusOK, "artist", data)
	})
}
