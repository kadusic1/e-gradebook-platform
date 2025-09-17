package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser     = flag.String("dbuser", "root", "Database user (default: root)")
	dbPassword = flag.String("dbpass", "", "Database password (leave empty for socket auth)")
	dbHost     = flag.String("dbhost", "127.0.0.1", "Database host (default: 127.0.0.1)")
	dbPort     = flag.String("dbport", "3306", "Database port (default: 3306)")
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

func getDBConfig() DBConfig {
	return DBConfig{
		User:     *dbUser,
		Password: *dbPassword,
		Host:     *dbHost,
		Port:     *dbPort,
	}
}

func (c DBConfig) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", c.User, c.Password, c.Host, c.Port)
}

func findEdnevnikDatabases(db *sql.DB) ([]string, error) {
	query := "SHOW DATABASES LIKE 'ednevnik_%'"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying databases: %v", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("error scanning database name: %v", err)
		}
		databases = append(databases, dbName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through database results: %v", err)
	}

	return databases, nil
}

func clearPreviousDumps() error {
	dumpsDir := "."
	files, err := filepath.Glob(filepath.Join(dumpsDir, "../ednevnik_*.sql"))
	if err != nil {
		return fmt.Errorf("error finding previous dump files: %v", err)
	}

	for _, file := range files {
		log.Printf("Removing previous dump file: %s", file)
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("error removing file %s: %v", file, err)
		}
	}

	return nil
}

func dumpDatabase(dbName string, config DBConfig) error {
	dumpsDir := ".."
	outputFile := filepath.Join(dumpsDir, fmt.Sprintf("%s.sql", dbName))

	log.Printf("Dumping database %s to %s", dbName, outputFile)

	// Build mysqldump command with comprehensive options
	args := []string{
		fmt.Sprintf("--user=%s", config.User),
		fmt.Sprintf("--password=%s", config.Password),
		fmt.Sprintf("--host=%s", config.Host),
		fmt.Sprintf("--port=%s", config.Port),
		"--single-transaction", // Consistent snapshot
		"--routines",           // Include stored procedures and functions
		"--triggers",           // Include triggers
		"--events",             // Include events
		"--create-options",     // Include all CREATE options
		"--extended-insert",    // Use extended INSERT syntax
		"--disable-keys",       // Disable keys during import
		"--lock-tables=false",  // Don't lock tables (use single-transaction instead)
		"--add-drop-database",  // Add DROP DATABASE statements
		"--databases",          // Dump specific database
		dbName,
	}

	cmd := exec.Command("mysqldump", args...)

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %v", outputFile, err)
	}
	defer outFile.Close()

	// Set command output to file
	cmd.Stdout = outFile

	// Capture stderr for error reporting
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Execute command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running mysqldump for %s: %v\nStderr: %s", dbName, err, stderr.String())
	}

	log.Printf("Successfully dumped database %s", dbName)
	return nil
}

func main() {
	flag.Parse()
	log.Println("Starting database dump process...")

	// Get database configuration
	config := getDBConfig()

	// Connect to database
	db, err := sql.Open("mysql", config.getDSN())
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Find all ednevnik databases
	databases, err := findEdnevnikDatabases(db)
	if err != nil {
		log.Fatalf("Error finding ednevnik databases: %v", err)
	}

	if len(databases) == 0 {
		log.Println("No ednevnik databases found")
		return
	}

	log.Printf("Found %d ednevnik databases: %v", len(databases), databases)

	// Clear previous dumps
	if err := clearPreviousDumps(); err != nil {
		log.Fatalf("Error clearing previous dumps: %v", err)
	}

	// Dump each database
	for _, dbName := range databases {
		if err := dumpDatabase(dbName, config); err != nil {
			log.Fatalf("Error dumping database %s: %v", dbName, err)
		}
	}

	log.Printf("Successfully completed dump process. Created %d dump files.", len(databases))
}
