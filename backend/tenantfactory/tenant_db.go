package tenantfactory

import (
	"database/sql"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
	"os"
)

// GetDBName TODO: Add description
func (t *ConfigurableTenant) GetDBName() (string, error) {
	safeName := util.SanitizeString(
		fmt.Sprintf("%d", t.TenantData.ID),
	)
	if safeName == "" {
		return "", fmt.Errorf("invalid tenant email for database name: %s", t.TenantData.Email)
	}
	return t.GetDBPrefix() + safeName, nil
}

// CreateSchema TODO: Add description
func (t *ConfigurableTenant) CreateSchema(tenantAdmin wpmodels.Teacher) error {
	content, err := os.ReadFile(t.Config.SchemaFile)
	if err != nil {
		return err
	}

	err = util.ExecSQLStatements(t.UserTenantDB, content)
	if err != nil {
		return fmt.Errorf("error executing SQL statements: %v", err)
	}

	dbName, err := t.GetDBName()
	if err != nil {
		return fmt.Errorf("error getting tenant DB name: %v", err)
	}

	// Grant select permissions to the service user
	err = util.GrantServiceReaderPrivileges(
		dbName, t.UserWorkspaceDB,
	)
	if err != nil {
		return fmt.Errorf("error granting service reader privileges: %v", err)
	}

	err = t.GrantTenantDBPrivileges()
	if err != nil {
		return fmt.Errorf("error granting tenant DB privileges: %v", err)
	}

	return nil
}

// CreateDB TODO: Add description
func (t *ConfigurableTenant) CreateDB() (string, error) {
	dbName, err := util.CreateTenantDB(t.GetDBPrefix(), t.TenantData.Email, t.UserWorkspaceDB)
	if err != nil {
		return "", err
	}

	return dbName, nil
}

// DropDB TODO: Add description
func (t *ConfigurableTenant) DropDB() error {
	err := t.RevokeTenantDBPrivileges()
	if err != nil {
		return fmt.Errorf("error revoking tenant DB privileges: %v", err)
	}

	dbName, err := t.GetDBName()
	if err != nil {
		return fmt.Errorf("error getting tenant DB name: %v", err)
	}

	// Revoke service user privileges
	err = util.RevokeServiceReaderPrivileges(dbName, t.UserWorkspaceDB)
	if err != nil {
		return fmt.Errorf("error revoking service reader privileges: %v", err)
	}

	return util.DropTenantDB(
		t.GetDBPrefix(),
		fmt.Sprintf("%d", t.TenantData.ID),
		t.UserWorkspaceDB,
	)
}

// GetDBPrefix TODO: Add description
func (t *ConfigurableTenant) GetDBPrefix() string {
	return t.Config.DBPrefix
}

// GrantTenantDBPrivileges TODO: Add description
func (t *ConfigurableTenant) GrantTenantDBPrivileges() error {
	tenantDBName, err := t.GetDBName()
	if err != nil {
		return fmt.Errorf("error getting tenant DB name: %v", err)
	}

	// Grant all privileges to tenant_admin
	var queries []string
	queries = append(queries, fmt.Sprintf(
		"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost' WITH GRANT OPTION;", tenantDBName, "tenant_admin",
	))

	// Grant privileges to teacher
	for _, p := range util.GetTeacherTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"GRANT %s ON %s.%s TO '%s'@'localhost' WITH GRANT OPTION;",
			p.Actions, tenantDBName, p.Table, "teacher",
		))
	}

	// Grant privileges to pupil
	for _, p := range util.GetPupilTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"GRANT %s ON %s.%s TO '%s'@'localhost';", p.Actions, tenantDBName, p.Table, "pupil",
		))
	}

	for _, q := range queries {
		if _, err := t.UserWorkspaceDB.Exec(q); err != nil {
			return fmt.Errorf("error executing query '%s': %v", q, err)
		}
	}
	return nil
}

// RevokeTenantDBPrivileges TODO: Add description
func (t *ConfigurableTenant) RevokeTenantDBPrivileges() error {
	var queries []string

	tenantDBName, err := t.GetDBName()
	if err != nil {
		return fmt.Errorf("error getting tenant DB name: %v", err)
	}

	// Revoke tenant admin privileges
	queries = append(queries, fmt.Sprintf(
		"REVOKE ALL PRIVILEGES ON %s.* FROM '%s'@'localhost';", tenantDBName, "tenant_admin",
	))

	// Revoke teacher privileges
	for _, p := range util.GetTeacherTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"REVOKE %s ON %s.%s FROM '%s'@'localhost';", p.Actions, tenantDBName, p.Table, "teacher",
		))
	}

	// Revoke pupil privileges
	for _, p := range util.GetPupilTablePrivileges() {
		queries = append(queries, fmt.Sprintf(
			"REVOKE %s ON %s.%s FROM '%s'@'localhost';", p.Actions, tenantDBName, p.Table, "pupil",
		))
	}

	for _, q := range queries {
		if _, err := t.UserWorkspaceDB.Exec(q); err != nil {
			return fmt.Errorf("error executing query '%s': %v", q, err)
		}
	}
	return nil
}

// StartTransactions begins database transactions on both the workspace and tenant databases.
// It returns the workspace transaction first, then the tenant transaction.
// If either transaction fails to start, it returns an error and both return values are nil.
//
// Returns:
//   - workspaceTx: Transaction for the workspace database
//   - tenantTx: Transaction for the tenant database
//   - error: Non-nil if either transaction failed to begin
//
// Note: Callers are responsible for properly committing or rolling back both transactions.
// Consider using defer statements to ensure rollback on error conditions.
func (t *ConfigurableTenant) StartTransactions(
	workspaceDb *sql.DB, tenantDb *sql.DB,
) (
	workspaceTx *sql.Tx, tenantTx *sql.Tx, err error,
) {
	tenantTx, err = tenantDb.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting tenant db transaction: %w", err)
	}

	workspaceTx, err = workspaceDb.Begin()
	if err != nil {
		tenantTx.Rollback() // Clean up the first transaction
		return nil, nil, fmt.Errorf("error starting workspace db transaction: %w", err)
	}

	return workspaceTx, tenantTx, nil
}
