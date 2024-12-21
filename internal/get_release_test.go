package internal

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetRelease(t *testing.T) {
	// Create a new sqlmock instance with exact matching for queries
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	// Create a mock logger
	e := echo.New()
	logger := e.Logger

	// Define the expected query and result for a valid release
	expectedQuery := "SELECT id, title, year FROM release WHERE id = $1;"
	mock.ExpectQuery(expectedQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "year"}).AddRow(1, "Bohemian Rhapsody", 1975))

	// Test valid release retrieval
	release, err := getRelease(db, 1, logger)
	assert.NoError(t, err)
	assert.NotNil(t, release)
	assert.Equal(t, 1, release.ID)
	assert.Equal(t, "Bohemian Rhapsody", release.Title)
	assert.Equal(t, "1975", release.Year)

	// Define the expected query for a non-existent release
	mock.ExpectQuery(expectedQuery).
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"})) // No rows

	// Test non-existent release
	release, err = getRelease(db, 999, logger)
	assert.Error(t, err)
	assert.Nil(t, release)
	assert.Contains(t, err.Error(), "release with ID 999 not found")

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
