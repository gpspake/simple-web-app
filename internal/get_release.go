package internal

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
)

// Release represents the release structure.
type Release struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// getRelease retrieves release information by ID from the database.
func getRelease(db DBQuerier, releaseID int, logger echo.Logger) (*Release, error) {
	// Validate input
	if releaseID <= 0 {
		return nil, fmt.Errorf("invalid release ID: %d", releaseID)
	}

	// Query to retrieve release information
	query := `
		SELECT id, title
		FROM release
		WHERE id = $1;
	`

	logger.Printf("Executing query: %s with releaseID: %d", query, releaseID)

	// Execute the query
	var release Release
	err := db.QueryRow(query, releaseID).Scan(&release.ID, &release.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("release with ID %d not found", releaseID)
		}
		return nil, fmt.Errorf("failed to retrieve release: %w", err)
	}

	return &release, nil
}
