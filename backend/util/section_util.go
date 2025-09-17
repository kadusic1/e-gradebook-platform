package util

import (
	"database/sql"
	"ednevnik-backend/config"
	commonmodels "ednevnik-backend/models/common"
	"ednevnik-backend/models/interfaces"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"strconv"
)

// GetSectionByID TODO: Add description
func GetSectionByID(
	sectionID int64,
	tenantDb interfaces.DatabaseQuerier,
) (tenantmodels.Section, error) {
	var section tenantmodels.Section
	query := `SELECT s.id, s.section_code, s.class_code, s.year, s.tenant_id,
	s.curriculum_code, t.id, CONCAT(t.name, ' ', t.last_name) as teacher_full_name,
	a.email, c.curriculum_name, s.archived, COALESCE(cs.course_name, '') as course_name
	FROM sections s
	JOIN ednevnik_workspace.curriculum c ON c.curriculum_code = s.curriculum_code
	LEFT JOIN homeroom_assignments ha ON ha.section_id = s.id
	LEFT JOIN ednevnik_workspace.teachers t ON t.id = ha.teacher_id
	LEFT JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	LEFT JOIN ednevnik_workspace.courses_secondary cs ON cs.course_code = c.course_code
	WHERE s.id = ?`

	var teacherID sql.NullInt64
	var teacherFullName sql.NullString
	var teacherEmail sql.NullString

	err := tenantDb.QueryRow(query, sectionID).Scan(
		&section.ID,
		&section.SectionCode,
		&section.ClassCode,
		&section.Year,
		&section.TenantID,
		&section.CurriculumCode,
		&teacherID,
		&teacherFullName,
		&teacherEmail,
		&section.CurriculumName,
		&section.Archived,
		&section.CourseName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return section, nil // No section found
		}
		return section, err // Other error
	}

	if teacherID.Valid {
		section.HomeroomTeacherID = int(teacherID.Int64)
	}
	if teacherFullName.Valid {
		section.HomeroomTeacherFullName = teacherFullName.String
	}
	if teacherEmail.Valid {
		section.HomeroomTeacherEmail = teacherEmail.String
	}

	// Populate additional fields
	section.Name = fmt.Sprintf("Odjeljenje %s-%s", section.ClassCode, section.SectionCode)
	return section, nil
}

func processSectionRows(
	rows *sql.Rows,
) ([]tenantmodels.Section, error) {
	var sections []tenantmodels.Section
	for rows.Next() {
		var section tenantmodels.Section
		var teacherID sql.NullInt64
		var teacherFullName sql.NullString
		var teacherEmail sql.NullString
		err := rows.Scan(
			&section.ID,
			&section.SectionCode,
			&section.ClassCode,
			&section.Year,
			&section.TenantID,
			&section.CurriculumCode,
			&teacherID,
			&teacherFullName,
			&teacherEmail,
			&section.CurriculumName,
		)
		if err != nil {
			return nil, err
		}
		if teacherID.Valid {
			section.HomeroomTeacherID = int(teacherID.Int64)
		}
		if teacherFullName.Valid {
			section.HomeroomTeacherFullName = teacherFullName.String
		}
		if teacherEmail.Valid {
			section.HomeroomTeacherEmail = teacherEmail.String
		}

		section.Name = fmt.Sprintf("Odjeljenje %s-%s", section.ClassCode, section.SectionCode)
		sections = append(sections, section)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}

// ListSectionsForTenant TODO: Add description
func ListSectionsForTenant(
	tenantID string, archived int, tenantDb *sql.DB, userWorkspaceDb *sql.DB,
) ([]tenantmodels.Section, error) {

	query := `SELECT s.id, s.section_code, s.class_code, s.year, s.tenant_id,
	s.curriculum_code, t.id, CONCAT(t.name, ' ', t.last_name) as teacher_full_name,
	a.email, c.curriculum_name
	FROM sections s
	JOIN ednevnik_workspace.curriculum c ON c.curriculum_code = s.curriculum_code
	LEFT JOIN homeroom_assignments ha ON ha.section_id = s.id
	LEFT JOIN ednevnik_workspace.teachers t ON t.id = ha.teacher_id
	LEFT JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	WHERE s.tenant_id = ? AND s.archived = ?`

	rows, err := tenantDb.Query(query, tenantID, archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections, err := processSectionRows(rows)
	if err != nil {
		return nil, err
	}

	return sections, err
}

// ListSectionsForTenantTeacher TODO: Add description
func ListSectionsForTenantTeacher(
	tenantID,
	teacherID,
	colorConfig string,
	archived int,
	tenantDb *sql.DB,
	userWorkspaceDb *sql.DB,
) ([]tenantmodels.Section, error) {
	query := `SELECT s.id, s.section_code, s.class_code, s.year, s.tenant_id,
	s.curriculum_code, t.id, CONCAT(t.name, ' ', t.last_name) as teacher_full_name,
	a.email, c.curriculum_name
	FROM sections s
	JOIN teachers_sections ts ON ts.section_id = s.id
	JOIN ednevnik_workspace.curriculum c ON c.curriculum_code = s.curriculum_code
	LEFT JOIN homeroom_assignments ha ON ha.section_id = s.id
	LEFT JOIN ednevnik_workspace.teachers t ON t.id = ha.teacher_id
	LEFT JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	WHERE s.tenant_id = ? AND ts.teacher_id = ? AND s.archived = ?`

	rows, err := tenantDb.Query(query, tenantID, teacherID, archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections, err := processSectionRows(rows)
	if err != nil {
		return nil, err
	}

	// Apply color config to all sections
	for i := range sections {
		sections[i].ColorConfig = colorConfig
	}

	return sections, nil
}

// ListSectionsForTenantPupil TODO: Add description
func ListSectionsForTenantPupil(
	tenantID,
	pupilID,
	colorConfig string,
	archived int,
	tenantDb *sql.DB,
	userWorkspaceDb *sql.DB,
) ([]tenantmodels.Section, error) {
	query := `SELECT s.id, s.section_code, s.class_code, s.year, s.tenant_id,
	s.curriculum_code, t.id, CONCAT(t.name, ' ', t.last_name) as teacher_full_name,
	a.email, c.curriculum_name
	FROM sections s
	JOIN pupils_sections ps ON ps.section_id = s.id
	JOIN ednevnik_workspace.curriculum c ON c.curriculum_code = s.curriculum_code
	LEFT JOIN homeroom_assignments ha ON ha.section_id = s.id
	LEFT JOIN ednevnik_workspace.teachers t ON t.id = ha.teacher_id
	LEFT JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	WHERE s.tenant_id = ? AND ps.pupil_id = ? AND s.archived = ?
	AND ps.is_active = 1 ORDER BY s.class_code ASC, s.section_code ASC`

	rows, err := tenantDb.Query(query, tenantID, pupilID, archived)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sections, err := processSectionRows(rows)
	if err != nil {
		return nil, err
	}

	// Apply color config to all sections
	for i := range sections {
		sections[i].ColorConfig = colorConfig
	}

	return sections, nil
}

// DeleteSection TODO: Add description
func DeleteSection(
	sectionID string, tenantDb *sql.DB,
) error {
	query := `DELETE FROM sections WHERE id = ?`
	_, err := tenantDb.Exec(query, sectionID)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}
	return nil
}

// UpdateSection TODO: Add description
func UpdateSection(
	section tenantmodels.Section,
	tenantDb *sql.DB,
	userWorkspaceDb *sql.DB,
	sectionID string,
) (tenantmodels.Section, error) {
	query := `UPDATE sections SET section_code = ?, class_code = ?, year = ?,
	curriculum_code = ?
	WHERE id = ?`

	_, err := tenantDb.Exec(query, section.SectionCode, section.ClassCode,
		section.Year, section.CurriculumCode, sectionID)
	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf("failed to update section: %w", err)
	}

	intSectionID, err := strconv.ParseInt(sectionID, 10, 64)
	if err != nil {
		return tenantmodels.Section{}, fmt.Errorf("invalid section_id: %w", err)
	}

	return GetSectionByID(
		intSectionID, tenantDb,
	)
}

// InviteAccountRecord TODO: Add description
type InviteAccountRecord struct {
	InviteID  int `json:"invite_id"`
	AccountID int `json:"account_id"`
}

// GetAllGlobalInviteIDsForSection TODO: Add description
func GetAllGlobalInviteIDsForSection(
	sectionID string,
	tenantDb *sql.DB,
	workspaceDB *sql.DB,
) ([]InviteAccountRecord, error) {
	var inviteIDs []InviteAccountRecord

	// First get all pupil invite IDs for the section
	pupilQuery := `SELECT psi.id, p.account_id FROM pupils_sections_invite psi
	JOIN pupils p ON p.id = psi.pupil_id
	WHERE section_id = ?`
	pupilRows, err := tenantDb.Query(pupilQuery, sectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invites for section: %w", err)
	}
	defer pupilRows.Close()

	for pupilRows.Next() {
		var inviteID int
		var accountID int
		if err := pupilRows.Scan(&inviteID, &accountID); err != nil {
			return nil, fmt.Errorf("failed to scan invite ID: %w", err)
		}
		inviteIDs = append(inviteIDs, InviteAccountRecord{
			InviteID:  inviteID,
			AccountID: accountID,
		})
	}
	if err := pupilRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over invites: %w", err)
	}

	// Now get all teacher invite IDs for the section
	teacherQuery := `SELECT ts.id, t.account_id FROM teachers_sections_invite ts
	JOIN ednevnik_workspace.teachers t ON t.id = ts.teacher_id
	WHERE section_id = ?`
	teacherRows, err := tenantDb.Query(teacherQuery, sectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher invites for section: %w", err)
	}
	defer teacherRows.Close()

	for teacherRows.Next() {
		var inviteID int
		var accountID int
		if err := teacherRows.Scan(&inviteID, &accountID); err != nil {
			return nil, fmt.Errorf("failed to scan teacher invite ID: %w", err)
		}
		inviteIDs = append(inviteIDs, InviteAccountRecord{
			InviteID:  inviteID,
			AccountID: accountID,
		})
	}
	if err := teacherRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over teacher invites: %w", err)
	}

	return inviteIDs, nil
}

// findFinalizedGradeHelper takes 3 arguments:
// pupilID - the ID of the pupil we are searching the final grade for
// subjectCode - the code of the subject we are searching the final grade for
// grades - a slice of final grades for a section
// If finalized grades are found the function returns true, else it returns
// false.
func findFinalizedGradeHelper(
	pupilID int, subjectCode, semesterCode string, grades []tenantmodels.Grade,
) bool {
	for _, grade := range grades {
		if grade.PupilID == pupilID && grade.SubjectCode == subjectCode && grade.SemesterCode == semesterCode {
			return true
		}
	}
	return false
}

// ArchiveSectionHelper updates the sections status to archived.
// It archives only if all pupils in the section have finalized grades for all
// section subjects.
// If the section is already archived it does nothing.
func ArchiveSectionHelper(
	sectionID int,
	tenant *wpmodels.Tenant,
	config *config.TenantConfig,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) error {
	var err error

	// Get section
	section, err := GetSectionByID(int64(sectionID), tenantDB)
	if err != nil {
		return err
	}

	// If section is archived do nothing
	if section.Archived {
		return nil
	}

	// Get section pupils
	pupils, err := GetPupilsForSection(
		fmt.Sprintf("%d", sectionID), false, tenantDB,
	)
	if err != nil {
		return err
	}
	if len(pupils) == 0 {
		return fmt.Errorf("prazno odjeljenje se ne može arhivirati - potrebno je da bude upisan bar jedan učenik")
	}

	// Get subjects for section
	subjects, err := GetAllSubjectsForCurriculumCode(
		section.CurriculumCode, workspaceDB,
	)
	if err != nil {
		return err
	}

	// Get semesters for section
	semesters, err := GetSemestersForSectionHelper(
		workspaceDB,
		tenantDB,
		fmt.Sprintf("%d", section.TenantID),
		fmt.Sprintf("%d", section.ID),
	)
	if err != nil {
		return err
	}

	// Get all final grades for section
	grades, err := GetAllFinalGradesForSection(sectionID, tenantDB)
	if err != nil {
		return err
	}

	// Initialize pupils without final grades just used for counting for now
	// but may be used in the future
	var pupilsWithoutFinalizedGrades []commonmodels.PupilSubjectSemester

	// Iterate over section pupils
	for _, pupil := range pupils {
		// Iterate over section subjects
		for _, subject := range subjects {
			for _, semester := range semesters {
				if !findFinalizedGradeHelper(
					pupil.ID, subject.SubjectCode, semester.SemesterCode, grades,
				) {
					pupilsWithoutFinalizedGrades = append(
						pupilsWithoutFinalizedGrades,
						commonmodels.PupilSubjectSemester{
							Pupil:    pupil,
							Subject:  subject,
							Semester: semester,
						},
					)
				}
			}
		}
	}

	if len(pupilsWithoutFinalizedGrades) > 0 {
		return fmt.Errorf(
			"odjeljenje se ne može arhivirati dok svi učenici nemaju zaključene ocjene iz svih predmeta",
		)
	}

	var finalCurricum bool
	finalCurriculumQuery := `SELECT final_curriculum FROM curriculum
	WHERE curriculum_code = ?`
	err = workspaceDB.QueryRow(
		finalCurriculumQuery, section.CurriculumCode,
	).Scan(&finalCurricum)
	if err != nil {
		return err
	}

	// Group grades by pupil_id for easier processing
	gradesByPupil := make(map[int][]tenantmodels.Grade)
	for _, grade := range grades {
		if grade.SemesterCode != config.MaxSemesterCode {
			continue
		}
		gradesByPupil[grade.PupilID] = append(gradesByPupil[grade.PupilID], grade)
	}

	// Get behaviour grades for section
	getBehaviourQuery := `SELECT id, pupil_id, section_id, behaviour, semester_code
	FROM pupil_behaviour WHERE section_id = ?`
	behaviourGrades := make(map[int]tenantmodels.BehaviourGrade)
	rows, err := tenantDB.Query(getBehaviourQuery, sectionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var grade tenantmodels.BehaviourGrade
		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID, &grade.Behaviour,
			&grade.SemesterCode,
		); err != nil {
			return err
		}
		if grade.SemesterCode == config.MaxSemesterCode {
			behaviourGrades[grade.PupilID] = grade
		}
	}

	// This query is used delete overlapping final grades if they exist for this
	// class and this school_specialization but this should not happen
	deleteFinalGradesGuardQuery := `DELETE FROM ` + config.FinalGradeTable +
		` WHERE pupil_id = ? AND class_code = ? AND school_specialization = ?`

	deleteBehaviourGradesGuardQuery := `DELETE FROM ` + config.BehaviourGradeTable + `
		WHERE pupil_id = ? AND class_code = ? AND school_specialization = ?`

	// Query to insert final grades
	insertFinalGradeQuery := `INSERT INTO ` + config.FinalGradeTable +
		` (pupil_id, tenant_id, subject_code, class_code, grade, school_specialization) 
		VALUES (?, ?, ?, ?, ?, ?)`

	// Query to insert behaviour grades
	insertBehaviourGradeQuery := `INSERT INTO ` + config.BehaviourGradeTable +
		` (pupil_id, tenant_id, class_code, behaviour, school_specialization) 
		VALUES (?, ?, ?, ?, ?)`

	updateEnrollmentStatusQuery := `UPDATE pupil_tenant SET ` +
		config.AvailableForEnrollmentField + ` = ? WHERE pupil_id = ?
		AND tenant_id = ?`

	// Start transaction
	workspaceTx, err := workspaceDB.Begin()
	defer func() {
		if err != nil {
			_ = workspaceTx.Rollback()
		}
	}()

	for _, pupil := range pupils {
		_, err = workspaceTx.Exec(
			deleteFinalGradesGuardQuery, pupil.ID, section.ClassCode,
			tenant.Specialization,
		)
		if err != nil {
			return err
		}

		_, err = workspaceTx.Exec(
			deleteBehaviourGradesGuardQuery, pupil.ID, section.ClassCode,
			tenant.Specialization,
		)
		if err != nil {
			return err
		}

		pupilGrades := gradesByPupil[pupil.ID]
		hasFailingGrade := false

		for _, grade := range pupilGrades {
			if grade.Grade < 2 {
				hasFailingGrade = true
				break
			}
		}

		if hasFailingGrade {
			continue
		}

		for _, grade := range pupilGrades {
			_, err = workspaceTx.Exec(
				insertFinalGradeQuery,
				grade.PupilID,
				tenant.ID,
				grade.SubjectCode,
				section.ClassCode,
				grade.Grade,
				tenant.Specialization,
			)
			if err != nil {
				return err
			}
		}

		behaviourGradeForPupil := behaviourGrades[pupil.ID]
		_, err = workspaceTx.Exec(
			insertBehaviourGradeQuery,
			behaviourGradeForPupil.PupilID,
			tenant.ID,
			section.ClassCode,
			behaviourGradeForPupil.Behaviour,
			tenant.Specialization,
		)
		if err != nil {
			return err
		}

		if finalCurricum {
			_, err = workspaceTx.Exec(
				updateEnrollmentStatusQuery,
				1,
				pupil.ID,
				tenant.ID,
			)
			if err != nil {
				return err
			}
		}
	}

	err = workspaceTx.Commit()
	if err != nil {
		return err
	}

	// Finally update section to be archived
	updateQuery := `UPDATE sections SET archived = ? WHERE id = ?`
	_, err = tenantDB.Exec(updateQuery, 1, sectionID)
	if err != nil {
		return err
	}

	return nil
}
