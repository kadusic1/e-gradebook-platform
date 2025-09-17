package tenantfactory

import (
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
)

// GetCurriculumsForAssignment TODO: Add description
func (t *ConfigurableTenant) GetCurriculumsForAssignment() ([]wpmodels.Curriculum, error) {
	var curriculums []wpmodels.Curriculum
	query := `SELECT c.curriculum_code, c.curriculum_name, c.class_code,
	c.npp_code, COALESCE(cs.course_name, ''), c.canton_code, c.tenant_type, n.npp_name
	FROM curriculum c
	JOIN npp n ON c.npp_code = n.npp_code
	LEFT JOIN courses_secondary cs ON c.course_code = cs.course_code
	WHERE c.tenant_type = ? AND NOT EXISTS (
		SELECT 1 FROM curriculum_tenant ct
		WHERE ct.curriculum_code = c.curriculum_code
		AND ct.tenant_id = ?
	)
	AND c.canton_code = ?`

	rows, err := t.UserWorkspaceDB.Query(
		query, t.TenantData.TenantType, t.TenantData.ID, t.TenantData.CantonCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get curriculums for assignment: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var curriculum wpmodels.Curriculum
		if err := rows.Scan(
			&curriculum.CurriculumCode,
			&curriculum.CurriculumName,
			&curriculum.ClassCode,
			&curriculum.NPPCode,
			&curriculum.CourseName,
			&curriculum.CantonCode,
			&curriculum.TenantType,
			&curriculum.NPPName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan curriculum: %w", err)
		}
		curriculums = append(curriculums, curriculum)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return curriculums, nil
}

// AssignCurriculumsToTenant TODO: Add description
func (t *ConfigurableTenant) AssignCurriculumsToTenant(curriculumCodes []string) error {
	if len(curriculumCodes) == 0 {
		return fmt.Errorf("no curriculum codes provided for assignment")
	}

	// Check if tenant with the given tenant_id exists
	var tenantCount int
	countQuery := `SELECT COUNT(*) FROM tenant WHERE id = ?`
	err := t.UserWorkspaceDB.QueryRow(countQuery, t.TenantData.ID).Scan(&tenantCount)
	if err != nil {
		return fmt.Errorf("failed to check tenant existence: %w", err)
	}
	if tenantCount == 0 {
		return fmt.Errorf("tenant with ID %d does not exist", t.TenantData.ID)
	}

	insertQuery := `INSERT INTO curriculum_tenant (curriculum_code, tenant_id)
	VALUES (?, ?)`

	for _, code := range curriculumCodes {
		_, err := t.UserWorkspaceDB.Exec(insertQuery, code, t.TenantData.ID)
		if err != nil {
			return fmt.Errorf("failed to assign curriculum %s to tenant %s: %w", code, t.TenantData.Email, err)
		}
	}

	err = util.TenantSemesterAssign(
		fmt.Sprintf("%d", t.TenantData.ID),
		t.UserWorkspaceDB,
	)
	if err != nil {
		return fmt.Errorf("failed to assign semesters to tenant: %v", err)
	}

	return nil
}

// GetCurriculumsForTenant TODO: Add description
func (t *ConfigurableTenant) GetCurriculumsForTenant() ([]wpmodels.CurriculumGet, error) {
	var curriculums []wpmodels.CurriculumGet
	query := `SELECT c.curriculum_code, c.curriculum_name, c.class_code, n.npp_name,
	cs.course_name, c.canton_code, c.tenant_type
	FROM curriculum c
	JOIN curriculum_tenant ct ON c.curriculum_code = ct.curriculum_code
	JOIN npp n ON c.npp_code = n.npp_code
	LEFT JOIN courses_secondary cs ON c.course_code = cs.course_code
	WHERE ct.tenant_id = ?`

	rows, err := t.UserWorkspaceDB.Query(query, t.TenantData.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get curriculums for tenant: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var curriculum wpmodels.CurriculumGet
		if err := rows.Scan(
			&curriculum.CurriculumCode,
			&curriculum.CurriculumName,
			&curriculum.ClassCode,
			&curriculum.NPPName,
			&curriculum.CourseName,
			&curriculum.CantonCode,
			&curriculum.TenantType,
		); err != nil {
			return nil, fmt.Errorf("failed to scan curriculum: %w", err)
		}
		curriculums = append(curriculums, curriculum)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return curriculums, nil
}

// UnassignCurriculumFromTenant TODO: Add description
func (t *ConfigurableTenant) UnassignCurriculumFromTenant(curriculumCode string) error {
	if curriculumCode == "" {
		return fmt.Errorf("missing curriculum code")
	}

	// Check if tenant with the given tenantID exists
	var tenantCount int
	countQuery := `SELECT COUNT(*) FROM tenant WHERE id = ?`
	err := t.UserWorkspaceDB.QueryRow(countQuery, t.TenantData.ID).Scan(&tenantCount)
	if err != nil {
		return fmt.Errorf("failed to check tenant existence: %w", err)
	}
	if tenantCount == 0 {
		return fmt.Errorf("tenant with ID %d does not exist", t.TenantData.ID)
	}

	// Check if section using this curriclum exists if so deny the deletion
	var sectionCount int
	sectionCheckQuery := `SELECT COUNT(*) FROM sections WHERE curriculum_code = ?`
	err = t.UserTenantDB.QueryRow(sectionCheckQuery, curriculumCode).Scan(&sectionCount)
	if err != nil {
		return fmt.Errorf("failed to check section existence: %v", err)
	}
	if sectionCount > 0 {
		return fmt.Errorf("kurikulum nije moguće izbrisati jer se trenutno koristi u jednom ili više odjeljenja")
	}

	deleteQuery := `DELETE FROM curriculum_tenant WHERE tenant_id = ? AND curriculum_code = ?`

	_, err = t.UserWorkspaceDB.Exec(deleteQuery, t.TenantData.ID, curriculumCode)
	if err != nil {
		return fmt.Errorf("failed to unassign curriculum %s from tenant %d: %w", curriculumCode, t.TenantData.ID, err)
	}

	err = util.TenantSemesterCleanup(
		fmt.Sprintf("%d", t.TenantData.ID),
		t.UserWorkspaceDB,
	)
	if err != nil {
		return fmt.Errorf("failed to cleanup tenant semesters: %v", err)
	}

	return nil
}
