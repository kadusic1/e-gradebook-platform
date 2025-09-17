package util

import (
	"database/sql"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
)

// GetGlobalDomainsHelper TODO: Add description
func GetGlobalDomainsHelper(workspaceDB *sql.DB) ([]wpmodels.Domain, error) {
	query := `SELECT domain FROM global_domains`
	rows, err := workspaceDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []wpmodels.Domain
	for rows.Next() {
		var domain wpmodels.Domain
		if err := rows.Scan(&domain.Domain); err != nil {
			return nil, err
		}
		domain.Type = "global_domain"
		domains = append(domains, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// GetTenantDomainsHelper TODO: Add description
func GetTenantDomainsHelper(workspaceDB *sql.DB) ([]wpmodels.Domain, error) {
	query := `SELECT domain FROM tenant WHERE domain IS NOT NULL
	AND domain <> ''`
	rows, err := workspaceDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []wpmodels.Domain
	for rows.Next() {
		var domain wpmodels.Domain
		if err := rows.Scan(&domain.Domain); err != nil {
			return nil, err
		}
		domain.Type = "tenant_domain"
		domains = append(domains, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return domains, nil
}

// GetAllDomainsHelper TODO: Add description
func GetAllDomainsHelper(workspaceDB *sql.DB) ([]wpmodels.Domain, error) {
	globalDomains, err := GetGlobalDomainsHelper(workspaceDB)
	if err != nil {
		return nil, err
	}

	tenantDomains, err := GetTenantDomainsHelper(workspaceDB)
	if err != nil {
		return nil, err
	}

	allDomains := append(globalDomains, tenantDomains...)
	return allDomains, nil
}

// InsertGlobalDomainHelper TODO: Add description
func InsertGlobalDomainHelper(workspaceDB *sql.DB, domain string) error {
	query := `SELECT COUNT(*) FROM global_domains WHERE domain = ?`
	var count int
	err := workspaceDB.QueryRow(query, domain).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("domena \"%s\" već postoji kao globalna domena", domain)
	}

	query = `SELECT COUNT(*) FROM tenant WHERE domain = ?`
	err = workspaceDB.QueryRow(query, domain).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("domena \"%s\" već postoji kao institucijska domena", domain)
	}

	query = `INSERT INTO global_domains (domain) VALUES (?)`
	_, err = workspaceDB.Exec(query, domain)
	return err
}

// DeleteGlobalDomainHelper TODO: Add description
func DeleteGlobalDomainHelper(workspaceDB *sql.DB, domain string) error {
	query := `DELETE FROM global_domains WHERE domain = ?`
	_, err := workspaceDB.Exec(query, domain)
	return err
}
