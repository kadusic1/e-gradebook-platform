package interfaces

import "database/sql"

// User TODO: Add description
type User interface {
	GetID() int
	GetName() string
	GetLastName() string
	GetEmail() string
	GetPhone() string
	GetAccountType(workspaceDB *sql.DB) string
	GetTenantIDs(workspaceDB *sql.DB) ([]string, error)
	GetPassword() string
	GetAccountID(workspaceDB *sql.DB) (int, error)
}
