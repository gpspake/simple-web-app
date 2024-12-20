package internal

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
)

// Artist represents the artist structure.
type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// getArtist retrieves artist information by ID from the database.
func getArtist(db DBQuerier, artistID int, logger echo.Logger) (*Artist, error) {
	// Validate input
	if artistID <= 0 {
		return nil, fmt.Errorf("invalid artist ID: %d", artistID)
	}

	// Query to retrieve artist information
	query := `
		SELECT id, name
		FROM artist
		WHERE id = $1;
	`

	logger.Printf("Executing query: %s with artistID: %d", query, artistID)

	// Execute the query
	var artist Artist
	err := db.QueryRow(query, artistID).Scan(&artist.ID, &artist.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artist with ID %d not found", artistID)
		}
		return nil, fmt.Errorf("failed to retrieve artist: %w", err)
	}

	return &artist, nil
}
