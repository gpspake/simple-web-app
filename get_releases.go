package main

import (
	"database/sql"
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
	// Default query to fetch releases with pagination
	query := "SELECT id, name, year FROM releases ORDER BY year ASC LIMIT ? OFFSET ?"
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		logger.Errorf("Error querying releases: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Collect results
	var items []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var year int
		err := rows.Scan(&id, &name, &year)
		if err != nil {
			logger.Errorf("Error scanning row: %v", err)
			return nil, err
		}

		items = append(items, map[string]interface{}{
			"id":   id,
			"name": name,
			"year": year,
		})
	}

	return items, nil
}
