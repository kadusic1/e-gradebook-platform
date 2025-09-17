package wpmodels

import (
	"database/sql"
	"ednevnik-backend/models/interfaces"
	"fmt"
)

// Teacher represents a teacher in the system
type Teacher struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password,omitempty"`
	Phone        string `json:"phone"`
	AccountType  string `json:"account_type,omitempty"`
	Contractions string `json:"contractions,omitempty"`
	Title        string `json:"title,omitempty"`
}

// GetID TODO: Add description
func (t Teacher) GetID() int {
	return t.ID
}

// GetName TODO: Add description
func (t Teacher) GetName() string {
	return t.Name
}

// GetLastName TODO: Add description
func (t Teacher) GetLastName() string {
	return t.LastName
}

// GetEmail TODO: Add description
func (t Teacher) GetEmail() string {
	return t.Email

}

// GetPhone TODO: Add description
func (t Teacher) GetPhone() string {
	return t.Phone
}

// GetAccountType TODO: Add description
func (t Teacher) GetAccountType(workspaceDB *sql.DB) string {
	var accountType string
	query := `SELECT account_type FROM accounts a
	JOIN teachers t ON a.id = t.account_id
	WHERE t.id = ?`

	err := workspaceDB.QueryRow(query, t.GetID()).Scan(&accountType)
	if err != nil {
		return ""
	}
	return accountType
}

// GetTenantIDs TODO: Add description
func (t Teacher) GetTenantIDs(workspaceDB *sql.DB) ([]string, error) {
	var tenantIDs []string

	query := "SELECT tenant_id FROM teacher_tenant WHERE teacher_id = ?"
	rows, err := workspaceDB.Query(query, t.GetID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tenantID string
		err := rows.Scan(&tenantID)
		if err != nil {
			return nil, err
		}
		tenantIDs = append(tenantIDs, tenantID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tenantIDs, nil
}

// GetPassword TODO: Add description
func (t Teacher) GetPassword() string {
	return t.Password
}

// GetAccountID TODO: Add description
func (t Teacher) GetAccountID(workspaceDB *sql.DB) (int, error) {
	var accountID int
	query := "SELECT account_id FROM teachers WHERE id = ?"
	err := workspaceDB.QueryRow(query, t.GetID()).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("error fetching account ID: %w", err)
	}
	return accountID, nil
}

// Ensure pupil implements user interface
var _ interfaces.User = (*Teacher)(nil)
