package interfaces

import "database/sql"

// DatabaseQuerier defines the interface for executing SQL queries.
// It abstracts the common query operations needed for database interactions,
// allowing for easy testing and dependency injection by accepting both
// *sql.DB and *sql.Tx implementations.
type DatabaseQuerier interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}
