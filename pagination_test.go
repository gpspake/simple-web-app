package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPagination(t *testing.T) {
	// Helper to create a mock HTTP request
	createRequest := func(queryParams map[string]string) *http.Request {
		q := url.Values{}
		for key, value := range queryParams {
			q.Set(key, value)
		}
		return &http.Request{
			URL: &url.URL{
				Path:     "/releases",
				RawQuery: q.Encode(),
			},
		}
	}

	t.Run("Custom Page and Limit", func(t *testing.T) {
		req := createRequest(map[string]string{"page": "2", "page_size": "20"})
		totalCount := 100
		pagination, err := getPagination("2", "20", totalCount, req)

		assert.NoError(t, err)
		assert.Equal(t, 2, pagination.Page)
		assert.Equal(t, 20, pagination.Limit)
		assert.Equal(t, 20, pagination.Offset)
		assert.Equal(t, 100, pagination.TotalCount)
		assert.Equal(t, 5, pagination.TotalPages)
		assert.Equal(t, "/releases?page=3&page_size=20", *pagination.NextUrl)
		assert.Equal(t, "/releases?page=1&page_size=20", *pagination.PrevUrl)
	})

	t.Run("First Page with Large Limit", func(t *testing.T) {
		req := createRequest(map[string]string{"page": "1", "page_size": "100"})
		totalCount := 50
		pagination, err := getPagination("1", "100", totalCount, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, pagination.Page)
		assert.Equal(t, 50, pagination.Limit) // Adjusted limit
		assert.Equal(t, 0, pagination.Offset)
		assert.Equal(t, 50, pagination.TotalCount)
		assert.Equal(t, 1, pagination.TotalPages)
		assert.Nil(t, pagination.NextUrl)
		assert.Nil(t, pagination.PrevUrl)
	})
}
