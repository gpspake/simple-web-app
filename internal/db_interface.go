// db_interface.go
package internal

import "database/sql"

// DBQuerier defines the subset of *sql.DB methods needed for querying
type DBQuerier interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
