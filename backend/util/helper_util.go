package util

import (
	"database/sql"
	"strings"
)

// IsDuplicateEmailError is helper function to check if error is duplicate email
func IsDuplicateEmailError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 1062") &&
		strings.Contains(errorMsg, "Duplicate entry") &&
		strings.Contains(errorMsg, "for key 'email'")
}

// AccountWithEmailExists TODO: Add description
func AccountWithEmailExists(email string, workspaceDB *sql.DB) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM accounts WHERE email = ?)"
	err := workspaceDB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// TeacherWithPhoneExists checks if a teacher with the given phone number exists
func TeacherWithPhoneExists(phone string, workspaceDB *sql.DB) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM teachers WHERE phone = ?)"
	err := workspaceDB.QueryRow(query, phone).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// PupilWithPhoneExists checks if a pupil with the given phone number exists
func PupilWithPhoneExists(phone string, workspaceDB *sql.DB) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pupil_global WHERE phone_number = ?)"
	err := workspaceDB.QueryRow(query, phone).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// IsDuplicatePhoneError checks if the error is due to a duplicate phone number
func IsDuplicatePhoneError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 1062") &&
		strings.Contains(errorMsg, "Duplicate entry") &&
		(strings.Contains(errorMsg, "for key 'phone'") ||
			strings.Contains(errorMsg, "for key 'phone_number'"))
}

// IsDuplicateDomain checks if the error is due to a duplicate domain
func IsDuplicateDomain(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 1062") &&
		strings.Contains(errorMsg, "Duplicate entry") &&
		strings.Contains(errorMsg, "for key 'domain'")
}

// GlobalDomainExists checks if a global domain exists in the database
func GlobalDomainExists(domain string, workspaceDB *sql.DB) (bool, error) {
	var exists bool
	if domain == "" {
		return false, nil
	}
	query := "SELECT EXISTS(SELECT 1 FROM global_domains WHERE domain = ?)"
	err := workspaceDB.QueryRow(query, domain).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// DuplicatePrimaryKeyHelper determines if an error is due to a duplicate PK
func DuplicatePrimaryKeyHelper(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 1062") &&
		strings.Contains(errorMsg, "Duplicate entry") &&
		strings.Contains(errorMsg, "for key 'PRIMARY'")
}

// InvalidSectionYearHelper checks if the constraint that checks that section
// format is YYYY/YYYY failed
func InvalidSectionYearHelper(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 4025") &&
		strings.Contains(errorMsg, "check_section_year")
}

// DuplicateSectionHelper checks for duplicate section DB constraint fail
func DuplicateSectionHelper(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "Error 1062") &&
		strings.Contains(errorMsg, "Duplicate entry") &&
		strings.Contains(errorMsg, "unique_section_class_year")
}
