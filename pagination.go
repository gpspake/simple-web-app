package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Pagination struct {
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
	TotalCount int     `json:"totalCount"`
	TotalPages int     `json:"totalPages"`
	First      int     `json:"first"`
	Last       int     `json:"last"`
	NextPage   *int    `json:"nextPage,omitempty"`
	PrevPage   *int    `json:"prevPage,omitempty"`
	NextUrl    *string `json:"nextUrl,omitempty"`
	PrevUrl    *string `json:"prevUrl,omitempty"`
}

func getPagination(
	pageStr string,
	limitStr string,
	totalCount int,
	request *http.Request,
) (Pagination, error) {
	const DefaultPageSize = 10

	// Parse inputs
	page := defaultInt(pageStr, 1)
	limit := defaultInt(limitStr, DefaultPageSize)

	if limit <= 0 {
		return Pagination{}, fmt.Errorf("invalid limit: %d", limit)
	}

	// Cap limit at totalCount to prevent excessive limits
	if limit > totalCount {

		// Prevent limit from being less than 1
		if totalCount < 1 {
			limit = 1
		} else {
			limit = totalCount
		}

	}

	// Calculate total pages and adjust page if out of bounds
	totalPages := (totalCount + limit - 1) / limit
	if page > totalPages {
		page = totalPages
	}
	if page < 1 {
		page = 1
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Calculate first and last item numbers
	first := offset + 1
	last := offset + limit
	if last > totalCount {
		last = totalCount
	}

	// Prepare URLs for next and previous pages
	var nextPage, prevPage *int
	var nextUrl, prevUrl *string

	queryParams := request.URL.Query()

	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal

		queryParams.Set("page", strconv.Itoa(*nextPage))
		nextUrlString := (&url.URL{
			Path:     request.URL.Path,
			RawQuery: queryParams.Encode(),
		}).String()
		nextUrl = &nextUrlString
	}

	if page > 1 {
		prevPageVal := page - 1
		prevPage = &prevPageVal

		queryParams.Set("page", strconv.Itoa(*prevPage))
		prevUrlString := (&url.URL{
			Path:     request.URL.Path,
			RawQuery: queryParams.Encode(),
		}).String()
		prevUrl = &prevUrlString
	}

	// Return pagination metadata
	return Pagination{
		Page:       page,
		Limit:      limit,
		Offset:     offset,
		TotalCount: totalCount,
		TotalPages: totalPages,
		First:      first,
		Last:       last,
		NextPage:   nextPage,
		PrevPage:   prevPage,
		NextUrl:    nextUrl,
		PrevUrl:    prevUrl,
	}, nil
}
