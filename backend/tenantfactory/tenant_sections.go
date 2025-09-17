package tenantfactory

import (
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
	"strings"
)

// GetSectionsForTenant TODO: Add description
func (t *ConfigurableTenant) GetSectionsForTenant(
	archived int,
) ([]tenantmodels.Section, error) {

	sections, err := util.ListSectionsForTenant(
		fmt.Sprintf("%d", t.TenantData.ID),
		archived,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sections for tenant: %v", err)
	}

	return sections, nil
}

// GetSectionsForTeacher TODO: Add description
func (t *ConfigurableTenant) GetSectionsForTeacher(
	teacherID string, archived int,
) ([]tenantmodels.Section, error) {

	sections, err := util.ListSectionsForTenantTeacher(
		fmt.Sprintf("%d", t.TenantData.ID),
		teacherID,
		t.TenantData.ColorConfig,
		archived,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sections for teacher %s: %v", teacherID, err)
	}
	return sections, nil
}

// GetSectionsForPupil TODO: Add description
func (t *ConfigurableTenant) GetSectionsForPupil(
	pupilID string, archived int,
) ([]tenantmodels.Section, error) {

	sections, err := util.ListSectionsForTenantPupil(
		fmt.Sprintf("%d", t.TenantData.ID),
		pupilID,
		t.TenantData.ColorConfig,
		archived,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting sections for pupil %s: %v", pupilID, err)
	}
	return sections, nil
}

// DeleteTenantSection TODO: Add description
func (t *ConfigurableTenant) DeleteTenantSection(sectionID string) error {
	// Get teacher in section for cleanup
	var teacherIDs []int
	teacherIDQuery := `SELECT teacher_id FROM teachers_sections
	WHERE section_id = ?`
	rows, err := t.UserTenantDB.Query(teacherIDQuery, sectionID)
	if err != nil {
		return fmt.Errorf("failed to get teachers for section %s: %w", sectionID, err)
	}
	defer rows.Close()
	for rows.Next() {
		var teacherID int
		if err := rows.Scan(&teacherID); err != nil {
			return fmt.Errorf("failed to scan teacher ID for section %s: %w", sectionID, err)
		}
		teacherIDs = append(teacherIDs, teacherID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over teachers for section %s: %w", sectionID, err)
	}

	// Get pupil IDs in section for cleanup
	pupilIDQuery := `SELECT pupil_id FROM pupils_sections
	WHERE section_id = ?`
	pupilRows, err := t.UserTenantDB.Query(pupilIDQuery, sectionID)
	if err != nil {
		return fmt.Errorf("failed to get pupils for section %s: %w", sectionID, err)
	}
	defer pupilRows.Close()
	var pupilIDs []int
	for pupilRows.Next() {
		var pupilID int
		if err := pupilRows.Scan(&pupilID); err != nil {
			return fmt.Errorf("failed to scan pupil ID for section %s: %w", sectionID, err)
		}
		pupilIDs = append(pupilIDs, pupilID)
	}

	inviteIndexIds, err := util.GetAllGlobalInviteIDsForSection(
		sectionID, t.UserTenantDB, t.UserWorkspaceDB,
	)
	if err != nil {
		return fmt.Errorf("failed to get global invite IDs for section %s: %w",
			sectionID, err,
		)
	}

	err = util.DeleteSection(sectionID, t.UserTenantDB)
	if err != nil {
		return fmt.Errorf("failed to delete section %s: %w", sectionID, err)
	}

	// Now for all teacher ids see if they are not assigned to any other section
	for _, teacherID := range teacherIDs {
		countQuery := `SELECT COUNT(*) FROM teachers_sections WHERE teacher_id = ?`
		var count int
		err = t.UserTenantDB.QueryRow(countQuery, teacherID).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to count sections for teacher %d: %w", teacherID, err)
		}
		if count == 0 {
			// If teacher is not assigned to any other section, delete the teacher
			// from the tenant
			deleteTeacherQuery := `DELETE FROM teacher_tenant WHERE teacher_id = ?
			AND tenant_id = ?`
			_, err = t.UserWorkspaceDB.Exec(deleteTeacherQuery, teacherID, t.TenantData.ID)
			if err != nil {
				return fmt.Errorf("failed to delete teacher from tenant %d: %w", teacherID, err)
			}
		}
	}

	// Now for all pupil ids see if they are not assigned to any other section
	for _, pupilID := range pupilIDs {
		countQuery := `SELECT COUNT(*) FROM pupils_sections WHERE pupil_id = ?`
		var count int
		err = t.UserTenantDB.QueryRow(countQuery, pupilID).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to count sections for pupil %d: %w", pupilID, err)
		}
		if count == 0 {
			// Delete pupil record from tenant DB
			deletePupilFromTenantDbQuery := `DELETE FROM pupils
			WHERE id = ?`
			_, err = t.UserTenantDB.Exec(deletePupilFromTenantDbQuery, pupilID)
			if err != nil {
				return fmt.Errorf("failed to delete pupil from tenant DB %d: %w", pupilID, err)
			}
			// If pupil is not assigned to any other section, delete the pupil
			// from the tenant
			deletePupilQuery := `DELETE FROM pupil_tenant WHERE pupil_id = ?
			AND tenant_id = ?`
			_, err = t.UserWorkspaceDB.Exec(deletePupilQuery, pupilID, t.TenantData.ID)
			if err != nil {
				return fmt.Errorf("failed to delete pupil from tenant %d: %w", pupilID, err)
			}
		}
	}

	if len(inviteIndexIds) > 0 {
		// Build placeholders for the IN clause
		placeholders := make([]string, len(inviteIndexIds))
		args := make([]interface{}, 0, len(inviteIndexIds)*3)

		for i, rec := range inviteIndexIds {
			placeholders[i] = "(?, ?, ?)"
			args = append(args, rec.InviteID, rec.AccountID, t.TenantData.ID)
		}

		deleteInviteQuery := fmt.Sprintf(`DELETE FROM invite_index 
        WHERE (invite_id, account_id, tenant_id) IN (%s)`,
			strings.Join(placeholders, ","))

		_, err := t.UserWorkspaceDB.Exec(deleteInviteQuery, args...)
		if err != nil {
			return fmt.Errorf("failed to delete global invites for section %s: %w", sectionID, err)
		}
	}

	return nil
}

// UpdateTenantSection TODO: Add description
func (t *ConfigurableTenant) UpdateTenantSection(
	newSection tenantmodels.Section, sectionID string,
) (tenantmodels.Section, error) {
	newSection, err := util.UpdateSection(
		newSection, t.UserTenantDB, t.UserWorkspaceDB, sectionID,
	)
	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf(
			"failed to update section %s: %w", sectionID, err,
		)
	}

	return newSection, nil
}

// CreateTenantSection TODO: Add description
func (t *ConfigurableTenant) CreateTenantSection(
	section tenantmodels.SectionCreate,
) (tenantmodels.Section, error) {

	sectionCreateQuery := `INSERT INTO sections (section_code, class_code, year,
	tenant_id, curriculum_code)
	VALUES (?, ?, ?, ?, ?)`

	res, err := t.UserTenantDB.Exec(sectionCreateQuery, section.SectionCode,
		section.ClassCode, section.Year, t.TenantData.ID,
		section.CurriculumCode)

	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf("failed to create section: %w", err)
	}

	sectionID, err := res.LastInsertId()
	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	createdSection, err := util.GetSectionByID(sectionID, t.UserTenantDB)
	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf("failed to get created section: %w", err)
	}

	return createdSection, nil
}

// GetMetadataForSectionCreation TODO: Add description
func (t *ConfigurableTenant) GetMetadataForSectionCreation() (*tenantmodels.SectionCreateMetadata, error) {
	classCodes := []string{"I", "II", "III", "IV"}

	// If tenant type is primary, append additional classes
	if t.TenantData.TenantType == "primary" {
		classCodes = append(classCodes, "V", "VI", "VII", "VIII", "IX")
	}

	// Convert to []wpmodels.Class
	var classes []wpmodels.Class
	for _, code := range classCodes {
		classes = append(classes, wpmodels.Class{ClassCode: code})
	}

	// Get teachers
	teachers, err := t.GetTeachersForTenant()
	if err != nil {
		return nil, fmt.Errorf("error getting teachers for tenant: %v", err)
	}

	// Get curriculums
	curriculums, err := t.GetCurriculumsForTenant()
	if err != nil {
		return nil, fmt.Errorf("error getting curriculums for tenant: %v", err)
	}

	var SectionCreateMetadata = tenantmodels.SectionCreateMetadata{
		Classes:     classes,
		Teachers:    teachers,
		Curriculums: curriculums,
	}

	return &SectionCreateMetadata, nil
}

// GetSubjectsForSection retrieves all subjects associated with a specific section ID
// from the tenant's database. Returns a slice of Subject or an error if the
// operation fails. It checks the user's claims to determine if they are a root or tenant admin
// or a teacher, and retrieves subjects accordingly.
func (t *ConfigurableTenant) GetSubjectsForSection(
	sectionID int, claims *wpmodels.Claims,
) ([]wpmodels.Subject, error) {
	section, err := util.GetSectionByID(int64(sectionID), t.UserTenantDB)
	if err != nil {
		return nil, err
	}

	var subjects []wpmodels.Subject

	if claims.AccountType == "root" || claims.AccountType == "tenant_admin" {
		subjects, err = util.GetAllSubjectsForCurriculumCode(
			section.CurriculumCode, t.UserWorkspaceDB,
		)
	} else {
		subjects, err = util.GetSectionSubjectsForTeacher(
			claims.ID, sectionID, t.UserTenantDB,
		)
	}

	if err != nil {
		return nil, err
	}

	return subjects, nil
}

// GetPupilCountForSection retrieves the count of pupils in a specific section
// identified by sectionID.
func (t *ConfigurableTenant) GetPupilCountForSection(sectionID int) (int, error) {
	countQuery := `SELECT COUNT(*) FROM pupils_sections WHERE section_id = ?`
	var count int
	err := t.UserTenantDB.QueryRow(countQuery, sectionID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pupils for section %d: %w", sectionID, err)
	}
	return count, nil
}

// ArchiveSection uses the ArchiveSectionHelper util function to archive a
// section
func (t *ConfigurableTenant) ArchiveSection(sectionID int) error {
	return util.ArchiveSectionHelper(
		sectionID,
		&t.TenantData,
		&t.Config,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
}

// UnenrollPupilFromSection unenrolls a pupil from a section by setting
// is_active field to 0
func (t *ConfigurableTenant) UnenrollPupilFromSection(pupilID, sectionID int) error {
	query := `UPDATE pupils_sections SET is_active = 0 WHERE pupil_id = ?
	AND section_id = ?`
	_, err := t.UserTenantDB.Exec(query, pupilID, sectionID)
	return err
}
