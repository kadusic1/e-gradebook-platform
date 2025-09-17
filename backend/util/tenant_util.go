package util

import (
	"database/sql"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
)

// GetTenantByID TODO: Add description
func GetTenantByID(tenantID string, db *sql.DB) (*wpmodels.Tenant, error) {
	var tenant wpmodels.Tenant
	query := `SELECT t.id, t.tenant_name, t.tenant_type, t.canton_code, t.address, t.phone,
	t.email, t.director_name, t.color_config, t.tenant_admin_id, t.lesson_display,
	t.absence_display, t.classroom_display, t.pupil_display, t.pupil_invite_display,
	c.canton_name, t.tenant_city, t.specialization
	FROM tenant t
	JOIN cantons c ON
	c.canton_code = t.canton_code
	WHERE t.id = ?`
	err := db.QueryRow(query, tenantID).Scan(
		&tenant.ID,
		&tenant.TenantName,
		&tenant.TenantType,
		&tenant.CantonCode,
		&tenant.Address,
		&tenant.Phone,
		&tenant.Email,
		&tenant.DirectorName,
		&tenant.ColorConfig,
		&tenant.TeacherID,
		&tenant.LessonDisplay,
		&tenant.AbsenceDisplay,
		&tenant.ClassroomDisplay,
		&tenant.PupilDisplay,
		&tenant.PupilInviteDisplay,
		&tenant.CantonName,
		&tenant.TenantCity,
		&tenant.Specialization,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tenant with ID %s not found", tenantID)
		}
		return nil, fmt.Errorf("error retrieving tenant: %v", err)
	}
	if tenant.TenantType == "osnovna škola" {
		tenant.TenantType = "primary"
	}
	if tenant.TenantType == "srednja škola" {
		tenant.TenantType = "secondary"
	}
	return &tenant, nil
}

// GetAllTenantIDs TODO: Add description
func GetAllTenantIDs(workspaceDb *sql.DB) ([]string, error) {
	query := `SELECT id FROM tenant`
	rows, err := workspaceDb.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tenant IDs: %v", err)
	}
	defer rows.Close()
	var tenantIDs []string
	for rows.Next() {
		var tenantID string
		if err := rows.Scan(&tenantID); err != nil {
			return nil, fmt.Errorf("error scanning tenant ID: %v", err)
		}
		tenantIDs = append(tenantIDs, tenantID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tenant IDs: %v", err)
	}
	return tenantIDs, nil
}
