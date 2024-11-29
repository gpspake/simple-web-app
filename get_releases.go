package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getReleasesCount(db *sql.DB) (int, error) {
	query := "SELECT COUNT(*) FROM releases"
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getPaginatedReleases(
	db *sql.DB,
	pageStr string,
	limitStr string,
	logger echo.Logger,
	request *http.Request,
) ([]map[string]interface{},
	Pagination,
	error,
) {
	totalCount, err := getReleasesCount(db)
	if err != nil {
		logger.Printf("Failed to get releases count: %v", err)
	}

	pagination, err := getPagination(
		pageStr,
		limitStr,
		totalCount,
		request,
	)

	releases, err := getReleases(db, pagination.Limit, pagination.Offset, logger)

	return releases, pagination, err
}

func getReleases(db *sql.DB, limit int, offset int, logger echo.Logger) ([]map[string]interface{}, error) {
	// Validate inputs
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit: %d", limit)
	}
	if offset < 0 {
		return nil, fmt.Errorf("invalid offset: %d", offset)
	}

	// Default query to fetch releases with pagination
	query := `
        SELECT
            releases.id AS release_id,
            releases.name AS release_name,
            releases.year AS release_year,
            artists.name AS artist_name
        FROM
            release_artists
        JOIN
            artists ON release_artists.artist_id = artists.id
        JOIN
            releases ON release_artists.release_id = releases.id
        ORDER BY release_year ASC 
        LIMIT ? 
        OFFSET ?;
    `
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		logger.Errorf("Error querying releases: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Collect results
	var items []map[string]interface{}
	for rows.Next() {
		var releaseId int
		var releaseName string
		var releaseYear int
		var artistName string
		err := rows.Scan(&releaseId, &releaseName, &releaseYear, &artistName)
		if err != nil {
			logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}

		items = append(items, map[string]interface{}{
			"releaseId":   releaseId,
			"releaseName": releaseName,
			"releaseYear": releaseYear,
			"artistName":  artistName,
		})
	}

	return items, nil
}
