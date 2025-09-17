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

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
}

type TablePrivilege struct {
	Table   string
	Actions string
}

var (
	dbUser     = flag.String("dbuser", "root", "Database user (default: root)")
	dbPassword = flag.String("dbpass", "", "Database password (leave empty for socket auth)")
	dbHost     = flag.String("dbhost", "127.0.0.1", "Database host (default: 127.0.0.1)")
	dbPort     = flag.String("dbport", "3306", "Database port (default: 3306)")
)

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

func getTeacherTablePrivileges() []TablePrivilege {
	return []TablePrivilege{
		{"pupils", "SELECT, INSERT, UPDATE, DELETE"},
		{"sections", "SELECT, UPDATE"},
		{"pupils_sections", "SELECT, INSERT, UPDATE, DELETE"},
		{"student_grades", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupils_sections_invite", "SELECT, INSERT, UPDATE, DELETE"},
		{"teachers_sections_invite", "SELECT"},
		{"teachers_sections_invite_subjects", "SELECT"},
		{"homeroom_assignments", "SELECT"},
		{"teachers_sections", "SELECT"},
		{"teachers_sections_subjects", "SELECT"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_attendance", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_behaviour", "SELECT, INSERT, UPDATE, DELETE"},
	}
}

func getPupilTablePrivileges() []TablePrivilege {
	return []TablePrivilege{
		{"pupils", "SELECT"},
		{"sections", "SELECT"},
		{"pupils_sections", "SELECT"},
		{"student_grades", "SELECT"},
		{"pupils_sections_invite", "SELECT"},
		{"homeroom_assignments", "SELECT"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT"},
		{"pupil_attendance", "SELECT"},
		{"pupil_behaviour", "SELECT"},
	}
}

func getServiceUserTablePrivileges() []TablePrivilege {
	return []TablePrivilege{
		{"pupils", "SELECT, INSERT, UPDATE, DELETE"},
		{"sections", "SELECT, UPDATE"},
		{"pupils_sections", "SELECT, INSERT, UPDATE, DELETE"},
		{"student_grades", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupils_sections_invite", "SELECT, INSERT, UPDATE, DELETE"},
		{"teachers_sections_invite", "SELECT, UPDATE"},
		{"teachers_sections_invite_subjects", "SELECT"},
		{"teachers_sections", "INSERT, SELECT"},
		{"teachers_sections_subjects", "INSERT, SELECT"},
		{"homeroom_assignments", "SELECT, INSERT, UPDATE"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_attendance", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_behaviour", "SELECT, INSERT, UPDATE, DELETE"},
	}
}

func findDumpFiles() ([]string, error) {
	files, err := filepath.Glob("../ednevnik_*.sql")
	if err != nil {
		return nil, fmt.Errorf("error finding dump files: %v", err)
	}
	return files, nil
}

func createUsers(db *sql.DB) error {
	log.Println("Creating database users and setting up privileges...")

	queries := []string{
		"SELECT '[LOG] Dropping user tenant_admin if exists...' AS info",
		"DROP USER IF EXISTS 'tenant_admin'@'localhost'",
		"SELECT '[LOG] Creating user tenant_admin...' AS info",
		"CREATE USER 'tenant_admin'@'localhost'",
		"SELECT '[LOG] Granting tenant admin privileges...' AS info",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.accounts TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.teachers TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.teacher_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.tenant_semester TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.curriculum_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pupil_global TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pupil_tenant TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_accounts TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_teachers TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.pending_pupil_global TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE, DELETE ON ednevnik_workspace.invite_index TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_final_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_behaviour_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_final_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_behaviour_grades TO 'tenant_admin'@'localhost' WITH GRANT OPTION",

		"SELECT '[LOG] Dropping user teacher if exists...' AS info",
		"DROP USER IF EXISTS 'teacher'@'localhost'",
		"SELECT '[LOG] Creating user teacher...' AS info",
		"CREATE USER 'teacher'@'localhost'",
		"SELECT '[LOG] Granting teacher privileges...' AS info",
		"GRANT DELETE ON ednevnik_workspace.pupil_tenant TO teacher WITH GRANT OPTION",
		"GRANT INSERT, DELETE ON ednevnik_workspace.invite_index TO teacher WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_final_grades TO 'teacher'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.primary_school_behaviour_grades TO 'teacher'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_final_grades TO 'teacher'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON ednevnik_workspace.high_school_behaviour_grades TO 'teacher'@'localhost' WITH GRANT OPTION",

		"SELECT '[LOG] Dropping user pupil if exists...' AS info",
		"DROP USER IF EXISTS 'pupil'@'localhost'",
		"SELECT '[LOG] Creating user pupil...' AS info",
		"CREATE USER 'pupil'@'localhost'",
		"SELECT '[LOG] Granting pupil privileges...' AS info",
		"GRANT SELECT on ednevnik_workspace.* TO 'pupil'@'localhost' WITH GRANT OPTION",

		"SELECT '[LOG] Dropping user service_reader if exists...' AS info",
		"DROP USER IF EXISTS 'service_reader'@'localhost'",
		"SELECT '[LOG] Creating user service_reader...' AS info",
		"CREATE USER 'service_reader'@'localhost'",
		"SELECT '[LOG] Granting service_reader privileges...' AS info",
		"GRANT TRIGGER ON *.* TO 'service_reader'@'localhost'",
		"SELECT '[LOG] Granting service DB user workspace privileges...' AS info",
		"GRANT INSERT ON ednevnik_workspace.pupil_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT ON ednevnik_workspace.teacher_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, UPDATE ON ednevnik_workspace.accounts TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT UPDATE ON ednevnik_workspace.pupil_global TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT UPDATE ON ednevnik_workspace.teachers TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT SELECT ON ednevnik_workspace.* TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT DELETE ON ednevnik_workspace.pupil_tenant TO 'service_reader'@'localhost' WITH GRANT OPTION",
		"GRANT INSERT, DELETE ON ednevnik_workspace.invite_index TO 'service_reader'@'localhost' WITH GRANT OPTION",

		"FLUSH PRIVILEGES",
	}

	for _, query := range queries {
		if strings.HasPrefix(query, "SELECT '[LOG]") {
			// Execute and display log messages
			rows, err := db.Query(query)
			if err != nil {
				log.Printf("Warning: Could not execute log query: %v", err)
				continue
			}
			var info string
			if rows.Next() {
				rows.Scan(&info)
				log.Println(info)
			}
			rows.Close()
		} else {
			// Execute other queries
			if _, err := db.Exec(query); err != nil {
				// Log warning but continue - some operations might fail if users don't exist, etc.
				log.Printf("Warning: Error executing query '%s': %v", query, err)
			}
		}
	}

	return nil
}

func grantTenantDBPrivileges(db *sql.DB, tenantDBs []string) error {
	log.Println("Granting privileges on tenant databases...")

	for _, tenantDB := range tenantDBs {
		log.Printf("Granting privileges for tenant database: %s", tenantDB)

		// Grant all privileges to tenant_admin
		query := fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO 'tenant_admin'@'localhost' WITH GRANT OPTION", tenantDB)
		if _, err := db.Exec(query); err != nil {
			log.Printf("Warning: Error granting tenant_admin privileges on %s: %v", tenantDB, err)
		}

		// Grant privileges to teacher
		for _, p := range getTeacherTablePrivileges() {
			query := fmt.Sprintf("GRANT %s ON %s.%s TO 'teacher'@'localhost' WITH GRANT OPTION",
				p.Actions, tenantDB, p.Table)
			if _, err := db.Exec(query); err != nil {
				log.Printf("Warning: Error granting teacher privileges on %s.%s: %v", tenantDB, p.Table, err)
			}
		}

		// Grant privileges to pupil
		for _, p := range getPupilTablePrivileges() {
			query := fmt.Sprintf("GRANT %s ON %s.%s TO 'pupil'@'localhost'",
				p.Actions, tenantDB, p.Table)
			if _, err := db.Exec(query); err != nil {
				log.Printf("Warning: Error granting pupil privileges on %s.%s: %v", tenantDB, p.Table, err)
			}
		}

		// Grant privileges to service_reader
		for _, p := range getServiceUserTablePrivileges() {
			query := fmt.Sprintf("GRANT %s ON %s.%s TO 'service_reader'@'localhost' WITH GRANT OPTION",
				p.Actions, tenantDB, p.Table)
			if _, err := db.Exec(query); err != nil {
				log.Printf("Warning: Error granting service_reader privileges on %s.%s: %v", tenantDB, p.Table, err)
			}
		}
	}

	// Flush privileges
	if _, err := db.Exec("FLUSH PRIVILEGES"); err != nil {
		return fmt.Errorf("error flushing privileges: %v", err)
	}

	return nil
}

func loadDumpFile(filename string, config DBConfig) error {
	log.Printf("Loading dump file: %s", filename)

	// Build mysql command
	args := []string{
		fmt.Sprintf("--user=%s", config.User),
		fmt.Sprintf("--host=%s", config.Host),
		fmt.Sprintf("--port=%s", config.Port),
		"--force", // Continue even if there are errors
	}

	if config.Password != "" {
		args = append(args, fmt.Sprintf("--password=%s", config.Password))
	}

	cmd := exec.Command("mysql", args...)

	// Read dump file
	dumpFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening dump file %s: %v", filename, err)
	}
	defer dumpFile.Close()

	// Set command input to dump file
	cmd.Stdin = dumpFile

	// Capture stderr for error reporting
	var stderr strings.Builder
	cmd.Stderr = &stderr

	// Execute command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error loading dump file %s: %v\nStderr: %s", filename, err, stderr.String())
	}

	log.Printf("Successfully loaded dump file: %s", filename)
	return nil
}

func extractTenantDBNames(dumpFiles []string) []string {
	var tenantDBs []string

	for _, file := range dumpFiles {
		// Remove .sql extension and extract database name
		dbName := strings.TrimSuffix(filepath.Base(file), ".sql")

		// Only include tenant databases (not workspace)
		if strings.HasPrefix(dbName, "ednevnik_tenant_db_") {
			tenantDBs = append(tenantDBs, dbName)
		}
	}

	return tenantDBs
}

func main() {
	flag.Parse()
	log.Println("Starting database load process...")

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

	// Find dump files
	dumpFiles, err := findDumpFiles()
	if err != nil {
		log.Fatalf("Error finding dump files: %v", err)
	}

	if len(dumpFiles) == 0 {
		log.Println("No dump files found")
		return
	}

	log.Printf("Found %d dump files: %v", len(dumpFiles), dumpFiles)

	// Separate workspace and tenant dump files
	var workspaceDump string
	var tenantDumps []string

	for _, dumpFile := range dumpFiles {
		dbName := strings.TrimSuffix(filepath.Base(dumpFile), ".sql")
		if dbName == "ednevnik_workspace" {
			workspaceDump = dumpFile
		} else {
			tenantDumps = append(tenantDumps, dumpFile)
		}
	}

	// Load workspace database first
	if workspaceDump != "" {
		log.Println("Loading workspace database first...")
		if err := loadDumpFile(workspaceDump, config); err != nil {
			log.Fatalf("Error loading workspace dump file %s: %v", workspaceDump, err)
		}
	} else {
		log.Println("Warning: No workspace dump file found (ednevnik_workspace.sql)")
	}

	// Load tenant databases
	log.Println("Loading tenant databases...")
	for _, dumpFile := range tenantDumps {
		if err := loadDumpFile(dumpFile, config); err != nil {
			log.Fatalf("Error loading tenant dump file %s: %v", dumpFile, err)
		}
	}

	log.Println("All dump files loaded successfully")

	// Create users and grant workspace privileges
	if err := createUsers(db); err != nil {
		log.Fatalf("Error creating users: %v", err)
	}

	// Extract tenant database names and grant privileges
	tenantDBs := extractTenantDBNames(dumpFiles)
	if len(tenantDBs) > 0 {
		if err := grantTenantDBPrivileges(db, tenantDBs); err != nil {
			log.Fatalf("Error granting tenant database privileges: %v", err)
		}
	}

	log.Printf("Successfully completed load process. Loaded %d dump files and set up user privileges.", len(dumpFiles))
}
