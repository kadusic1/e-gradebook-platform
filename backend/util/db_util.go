package util

import (
	"database/sql"
	"fmt"
	"strings"
)

// ExecSQLStatements executes an SQL file
func ExecSQLStatements(db *sql.DB, content []byte) error {
	// 1) Normalize newlines (handles Windows CRLF safely)
	src := strings.ReplaceAll(string(content), "\r\n", "\n")

	lines := strings.Split(src, "\n")
	var statements []string
	var currentStmt strings.Builder
	delimiter := ";"

	for _, rawLine := range lines {
		lineTrim := strings.TrimSpace(rawLine)

		// Skip comments / empty
		if lineTrim == "" || strings.HasPrefix(lineTrim, "--") || strings.HasPrefix(lineTrim, "#") {
			continue
		}

		// Handle DELIMITER directive (client-side)
		up := strings.ToUpper(lineTrim)
		if strings.HasPrefix(up, "DELIMITER ") {
			// set new delimiter (everything after the first space)
			delimiter = strings.TrimSpace(lineTrim[len("DELIMITER "):])
			continue
		}

		// Accumulate raw line exactly as-is (no trimming) to preserve inner semicolons
		currentStmt.WriteString(rawLine)
		currentStmt.WriteString("\n")

		// End-of-statement if the *trimmed* line ends with the delimiter OR equals it
		if strings.HasSuffix(lineTrim, delimiter) || lineTrim == delimiter {
			// Build the full statement
			stmt := currentStmt.String()

			// Trim surrounding whitespace first, then strip trailing delimiter (which may otherwise be preceded by \n)
			stmt = strings.TrimSpace(stmt)
			stmt = strings.TrimSuffix(stmt, delimiter)
			stmt = strings.TrimSpace(stmt)

			if stmt != "" {
				statements = append(statements, stmt)
			}
			currentStmt.Reset()
		}
	}

	// If anything remains without a trailing delimiter, flush it
	rest := strings.TrimSpace(currentStmt.String())
	if rest != "" {
		statements = append(statements, rest)
	}

	// Execute
	for i, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("error executing statement %d: %v\nSQL: %s", i+1, err, stmt)
		}
	}
	return nil
}

// CreateTenantDB TODO: Add description
func CreateTenantDB(prefix string, id string, workspaceDB *sql.DB) (string, error) {
	safeName := strings.ReplaceAll(strings.ReplaceAll(id, "@", "_"), ".", "_")
	dbName := prefix + safeName
	createDBSQL := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_bin",
		dbName,
	)
	_, err := workspaceDB.Exec(createDBSQL)
	if err != nil {
		return "", err
	}

	return dbName, nil
}

// DropTenantDB TODO: Add description
func DropTenantDB(prefix string, id string, workspaceDB *sql.DB) error {
	dbName := prefix + SanitizeString(id)
	_, err := workspaceDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
	return err
}

// ConnectToTenantDB TODO: Add description
func ConnectToTenantDB(dbName string) (*sql.DB, error) {
	tenantcs := BuildDBConnectionString(dbName)
	tenantDB, err := sql.Open("mysql", tenantcs)
	if err != nil {
		return nil, err
	}
	if err = tenantDB.Ping(); err != nil {
		tenantDB.Close()
		return nil, err
	}
	return tenantDB, nil
}

// SanitizeString TODO: Add description
func SanitizeString(email string) string {
	return strings.ReplaceAll(strings.ReplaceAll(email, "@", "_"), ".", "_")
}

// GrantServiceReaderPrivileges TODO: Add description
func GrantServiceReaderPrivileges(tenantDBName string, workspaceDB *sql.DB) error {
	var queries []string

	for _, p := range GetServiceUserTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"GRANT %s ON %s.%s TO '%s'@'localhost' WITH GRANT OPTION;", p.Actions, tenantDBName, p.Table, "service_reader",
		))
	}

	for _, q := range queries {
		if _, err := workspaceDB.Exec(q); err != nil {
			return fmt.Errorf("error executing query '%s': %v", q, err)
		}
	}

	return nil
}

// RevokeServiceReaderPrivileges TODO: Add description
func RevokeServiceReaderPrivileges(tenantDBName string, workspaceDB *sql.DB) error {
	var queries []string

	for _, p := range GetServiceUserTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"REVOKE %s ON %s.%s FROM '%s'@'localhost';", p.Actions, tenantDBName, p.Table, "service_reader",
		))
	}

	for _, q := range queries {
		if _, err := workspaceDB.Exec(q); err != nil {
			return fmt.Errorf("error executing query '%s': %v", q, err)
		}
	}

	return nil
}
