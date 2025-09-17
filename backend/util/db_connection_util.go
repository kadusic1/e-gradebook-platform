package util

import (
	"database/sql"
	"ednevnik-backend/constants"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// BuildDBConnectionString constructs a MARIADB connection string from
// environment variables and a given db name.
func BuildDBConnectionString(dbname string) string {
	user := os.Getenv("MARIADB_USER")
	password := os.Getenv("MARIADB_PASSWORD")
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)
}

// BuildDBConnectionStringWithUser TODO: Add description
func BuildDBConnectionStringWithUser(dbname, accountType string) string {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	// If account type is root use eacon user
	if accountType == "root" {
		accountType = "eacon"
		// Get password
		password := os.Getenv("MARIADB_PASSWORD")
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s", accountType, password, host, port, dbname,
		)
	}

	return fmt.Sprintf(
		"%s:@tcp(%s:%s)/%s", accountType, host, port, dbname,
	)
}

// BuildServiceReaderConnectionString TODO: Add description
func BuildServiceReaderConnectionString(dbname string) string {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	return fmt.Sprintf(
		"%s:@tcp(%s:%s)/%s", "service_reader", host, port, dbname,
	)
}

// GetUserWorkspaceDBFromContext TODO: Add description
func GetUserWorkspaceDBFromContext(r *http.Request) (*sql.DB, bool) {
	db, ok := r.Context().Value(constants.UserWorkspaceDBKey).(*sql.DB)
	return db, ok
}

// dBCache used to cache database connections
var dbCache sync.Map

func getOrCreateDBConnection(connectionString string, errorContext string) (*sql.DB, error) {
	if dbCachedConnection, ok := dbCache.Load(connectionString); ok {
		return dbCachedConnection.(*sql.DB), nil
	}

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close() // Clean up the connection
		return nil, fmt.Errorf("failed to ping database %s: %w", errorContext, err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)
	dbCache.Store(connectionString, db)
	return db, nil
}

// GetOrCreateDBConnection TODO: Add description
func GetOrCreateDBConnection(dbname, accountType string) (*sql.DB, error) {
	connectionString := BuildDBConnectionStringWithUser(dbname, accountType)
	errorContext := fmt.Sprintf("account type %s", accountType)
	return getOrCreateDBConnection(connectionString, errorContext)
}

// GetOrCreateDBConnectionServiceReader TODO: Add description
func GetOrCreateDBConnectionServiceReader(dbname string) (*sql.DB, error) {
	connectionString := BuildServiceReaderConnectionString(dbname)
	errorContext := "with service reader"
	return getOrCreateDBConnection(connectionString, errorContext)
}
