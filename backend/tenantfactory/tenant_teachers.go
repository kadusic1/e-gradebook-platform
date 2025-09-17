package tenantfactory

import (
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
)

// GetTeachersForTenant TODO: Add description
func (t *ConfigurableTenant) GetTeachersForTenant() ([]wpmodels.Teacher, error) {
	query := `SELECT t.id, t.name, t.last_name, a.email, t.phone, a.account_type as role,
		t.contractions, t.title
	    FROM teachers t
		JOIN accounts a ON a.id = t.account_id
        JOIN teacher_tenant tt ON t.id = tt.teacher_id
        WHERE tt.tenant_id = ?`

	rows, err := t.UserWorkspaceDB.Query(query, t.TenantData.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting teachers for tenant: %v", err)
	}
	defer rows.Close()

	teachers := []wpmodels.Teacher{}

	for rows.Next() {
		var t wpmodels.Teacher
		err := rows.Scan(
			&t.ID, &t.Name, &t.LastName, &t.Email, &t.Phone, &t.AccountType,
			&t.Contractions, &t.Title,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading teacher for tenant: %v", err)
		}
		teachers = append(teachers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return teachers, nil
}

// DeleteTenantTeacherData TODO: Add description
func (t *ConfigurableTenant) DeleteTenantTeacherData(teacherID string) error {
	// Start a transaction for data consistency
	tx, err := t.UserTenantDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() succeeds

	// 1. Delete homeroom assignments
	homeroomQuery := `DELETE FROM homeroom_assignments WHERE teacher_id = ?`
	_, err = tx.Exec(homeroomQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete homeroom assignments: %w", err)
	}

	// No need to delete section invites here, as we want to keep them for history
	// 2. Delete teacher section invites (will cascade to teachers_sections_invite_subjects)
	// inviteQuery := `DELETE FROM teachers_sections_invite WHERE teacher_id = ?`
	// _, err = tx.Exec(inviteQuery, teacherID)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete teacher section invites: %w", err)
	// }

	// 3. Delete from teachers_sections_subjects first (due to foreign key)
	sectionSubjectsQuery := `DELETE FROM teachers_sections_subjects WHERE teacher_id = ?`
	_, err = tx.Exec(sectionSubjectsQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher section subjects: %w", err)
	}

	// 4. Delete from teachers_sections
	sectionsQuery := `DELETE FROM teachers_sections WHERE teacher_id = ?`
	_, err = tx.Exec(sectionsQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher sections: %w", err)
	}

	// 5. Update student_grades to set teacher_id to NULL
	gradesQuery := `UPDATE student_grades SET teacher_id = NULL WHERE teacher_id = ?`
	_, err = tx.Exec(gradesQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to update student grades: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 7. Delete tenant teacher record from workspaceDB
	tenantTeacherQuery := `DELETE FROM teacher_tenant WHERE
	teacher_id = ? AND tenant_id = ?`
	_, err = t.UserWorkspaceDB.Exec(tenantTeacherQuery, teacherID, t.TenantData.ID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher tenant record: %w", err)
	}

	return nil
}

// DeleteTenantTeacherDataWithoutInvites TODO: Add description
func (t *ConfigurableTenant) DeleteTenantTeacherDataWithoutInvites(teacherID string) error {
	// Start a transaction for data consistency
	tx, err := t.UserTenantDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() succeeds

	// 1. Delete homeroom assignments
	homeroomQuery := `DELETE FROM homeroom_assignments WHERE teacher_id = ?`
	_, err = tx.Exec(homeroomQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete homeroom assignments: %w", err)
	}

	// 2. Delete from teachers_sections_subjects first (due to foreign key)
	sectionSubjectsQuery := `DELETE FROM teachers_sections_subjects WHERE teacher_id = ?`
	_, err = tx.Exec(sectionSubjectsQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher section subjects: %w", err)
	}

	// 3. Delete from teachers_sections
	sectionsQuery := `DELETE FROM teachers_sections WHERE teacher_id = ?`
	_, err = tx.Exec(sectionsQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher sections: %w", err)
	}

	// 5. Update student_grades to set teacher_id to NULL
	gradesQuery := `UPDATE student_grades SET teacher_id = NULL WHERE teacher_id = ?`
	_, err = tx.Exec(gradesQuery, teacherID)
	if err != nil {
		return fmt.Errorf("failed to update student grades: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 7. Delete tenant teacher record from workspaceDB
	tenantTeacherQuery := `DELETE FROM teacher_tenant WHERE
	teacher_id = ? AND tenant_id = ?`
	_, err = t.UserWorkspaceDB.Exec(tenantTeacherQuery, teacherID, t.TenantData.ID)
	if err != nil {
		return fmt.Errorf("failed to delete teacher tenant record: %w", err)
	}

	return nil
}

// DeleteTeacherFromTenant TODO: Add description
func (t *ConfigurableTenant) DeleteTeacherFromTenant(teacherID string) error {
	err := t.DeleteTenantTeacherDataWithoutInvites(teacherID)
	if err != nil {
		return fmt.Errorf("error deleting teacher tenant data: %v", err)
	}

	return nil
}
