package internal

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetArtist(t *testing.T) {
	// Create a new sqlmock instance with exact matching for queries
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	// Create a mock logger
	e := echo.New()
	logger := e.Logger

	// Define the expected query and result for a valid artist
	expectedQuery := `
		SELECT id, name
		FROM artist
		WHERE id = $1;
	`
	mock.ExpectQuery(expectedQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Freddie Mercury"))

	// Test valid artist retrieval
	artist, err := getArtist(db, 1, logger)
	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, 1, artist.ID)
	assert.Equal(t, "Freddie Mercury", artist.Name)

	// Define the expected query for a non-existent artist
	mock.ExpectQuery(expectedQuery).
		WithArgs(999).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"})) // No rows

	// Test non-existent artist
	artist, err = getArtist(db, 999, logger)
	assert.Error(t, err)
	assert.Nil(t, artist)
	assert.Contains(t, err.Error(), "artist with ID 999 not found")

	// Verify that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
