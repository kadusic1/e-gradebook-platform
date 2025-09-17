package util

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"strconv"
)

// GetPupilsForTenantSectionAssignment TODO: Add description
func GetPupilsForTenantSectionAssignment(
	sectionID string, tenantDB *sql.DB,
) ([]tenantmodels.Pupil, error) {
	pupils := []tenantmodels.Pupil{}

	query := `SELECT DISTINCT p.id, p.name, p.last_name, p.gender, p.address,
	p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
	p.religion, a.email, p.place_of_birth
	FROM ednevnik_workspace.pupil_global p
	JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
	WHERE NOT EXISTS (
		SELECT 1 FROM pupils_sections ps
		WHERE ps.pupil_id = p.id AND ps.section_id = ? AND ps.is_active = 1
	)
	AND NOT EXISTS (
		SELECT 1 FROM pupils_sections_invite psi
		WHERE psi.section_id = ? AND psi.pupil_id = p.id AND psi.status = 'pending'
	)
	AND (
		-- Either pupil is not in any section with same class code in same tenant
		NOT EXISTS (
			SELECT 1 FROM pupils_sections ps2
			JOIN sections s2 ON ps2.section_id = s2.id
			JOIN sections s_current ON s_current.id = ?
			WHERE ps2.pupil_id = p.id
			AND ps2.is_active = 1
			AND s2.class_code = s_current.class_code
			AND s2.tenant_id = s_current.tenant_id
		)
		OR
		-- Or pupil has failed (has at least one final grade < 2 in max progress level semester)
		EXISTS (
			SELECT 1 FROM student_grades sg
			JOIN sections s3 ON sg.section_id = s3.id
			JOIN sections s_current2 ON s_current2.id = ?
			JOIN ednevnik_workspace.semester sem ON sg.semester_code = sem.semester_code
			WHERE sg.pupil_id = p.id
			AND sg.type = 'final'
			AND sg.grade < 2
			AND s3.class_code = s_current2.class_code
			AND s3.tenant_id = s_current2.tenant_id
			AND sem.progress_level = (
				SELECT MAX(sem2.progress_level)
				FROM ednevnik_workspace.semester sem2
			)
		)
	)
	ORDER BY p.last_name, p.name`

	rows, err := tenantDB.Query(query, sectionID, sectionID, sectionID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying pupils for section assignment: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pupil tenantmodels.Pupil
		if err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.LastName,
			&pupil.Gender,
			&pupil.Address,
			&pupil.GuardianName,
			&pupil.PhoneNumber,
			&pupil.GuardianNumber,
			&pupil.DateOfBirth,
			&pupil.Religion,
			&pupil.Email,
			&pupil.PlaceOfBirth,
		); err != nil {
			return nil, fmt.Errorf("error scanning pupil row: %v", err)
		}
		pupils = append(pupils, pupil)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over pupil rows: %v", err)
	}

	return pupils, nil
}

// GetSectionPupilsWithPendingInvites TODO: Add description
func GetSectionPupilsWithPendingInvites(
	sectionID string, tenantDB *sql.DB,
) ([]tenantmodels.Pupil, error) {
	pupils := []tenantmodels.Pupil{}

	query := `SELECT DISTINCT p.id, p.name, p.last_name, p.gender, p.address,
	p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
	p.religion, a.email, p.place_of_birth
	FROM ednevnik_workspace.pupil_global p
	JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
	JOIN pupils_sections_invite psi ON p.id = psi.pupil_id
	WHERE psi.section_id = ? AND psi.status = 'pending'
	ORDER BY p.last_name, p.name`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying pupils with pending invites for section: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pupil tenantmodels.Pupil
		if err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.LastName,
			&pupil.Gender,
			&pupil.Address,
			&pupil.GuardianName,
			&pupil.PhoneNumber,
			&pupil.GuardianNumber,
			&pupil.DateOfBirth,
			&pupil.Religion,
			&pupil.Email,
			&pupil.PlaceOfBirth,
		); err != nil {
			return nil, fmt.Errorf("error scanning pupil row: %v", err)
		}
		pupils = append(pupils, pupil)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over pupil rows: %v", err)
	}

	return pupils, nil
}

// SendPupilSectionInvite TODO: Add description
func SendPupilSectionInvite(
	pupilID int, sectionID, tenantID string, tenantDB *sql.DB,
	workspaceDB *sql.DB, tenantName string,
) (*tenantmodels.PupilSectionInvite, error) {
	query := `INSERT INTO pupils_sections_invite (pupil_id, section_id)
	VALUES (?, ?)`

	res, err := tenantDB.Exec(query, pupilID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error sending pupil section invite: %v", err)
	}

	// Get the last inserted ID to use in invite_index
	inviteID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	var pupilAccountID int
	accountQuery := `SELECT account_id FROM pupil_global WHERE id = ?`
	err = workspaceDB.QueryRow(accountQuery, pupilID).Scan(&pupilAccountID)
	if err != nil {
		return nil, fmt.Errorf("error getting pupil account ID: %v", err)
	}

	inviteIndexQuery := `INSERT INTO invite_index (invite_id, account_id, tenant_id)
	VALUES (?, ?, ?)`

	_, err = workspaceDB.Exec(inviteIndexQuery, inviteID, pupilAccountID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("error inserting into invite index: %v", err)
	}

	tenantIDint, err := strconv.Atoi(tenantID)
	if err != nil {
		return nil, fmt.Errorf("error converting tenantID to int: %v", err)
	}

	newInvite, err := GetPupilSectionInvite(
		int(inviteID), tenantName, tenantDB, tenantIDint,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting new pupil section invite: %v", err)
	}

	return newInvite, nil
}

// GetPupilSectionInvite TODO: Add description
func GetPupilSectionInvite(
	inviteID int, tenantName string, tenantDB *sql.DB, tenantID int,
) (*tenantmodels.PupilSectionInvite, error) {
	getQuery := `SELECT psi.id, psi.pupil_id, psi.section_id, psi.invite_date,
	psi.status, CONCAT(p.name, ' ', p.last_name) AS pupil_full_name,
	CONCAT('Odjeljenje ', s.class_code, '-', s.section_code) AS section_name,
	a.email
	FROM pupils_sections_invite psi
	JOIN ednevnik_workspace.pupil_global p ON psi.pupil_id = p.id
	JOIN ednevnik_workspace.accounts a ON a.id = p.account_id
	JOIN sections s ON psi.section_id = s.id
	WHERE psi.id = ?`

	var invite tenantmodels.PupilSectionInvite
	err := tenantDB.QueryRow(getQuery, inviteID).Scan(
		&invite.ID,
		&invite.PupilID,
		&invite.SectionID,
		&invite.InviteDate,
		&invite.Status,
		&invite.PupilFullName,
		&invite.SectionName,
		&invite.PupilEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invite not found for ID: %d", inviteID)
		}
		return nil, fmt.Errorf("error scanning pupil section invite row: %v", err)
	}
	invite.TenantID = tenantID
	invite.TenantName = tenantName

	return &invite, nil
}

// GetPupilSectionInvitesForSectionHelper TODO: Add description
func GetPupilSectionInvitesForSectionHelper(
	sectionID string, tenantName string, tenantDB *sql.DB, tenantID int,
) ([]tenantmodels.PupilSectionInvite, error) {
	getQuery := `SELECT psi.id, psi.pupil_id, psi.section_id, psi.invite_date,
	psi.status, CONCAT(p.name, ' ', p.last_name) AS pupil_full_name,
	CONCAT('Odjeljenje ', s.class_code, '-', s.section_code) AS section_name,
	a.email
	FROM pupils_sections_invite psi
	JOIN ednevnik_workspace.pupil_global p ON psi.pupil_id = p.id
	JOIN ednevnik_workspace.accounts a ON a.id = p.account_id
	JOIN sections s ON psi.section_id = s.id
	WHERE psi.section_id = ?`

	var invites []tenantmodels.PupilSectionInvite

	rows, err := tenantDB.Query(getQuery, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying invites for sectionID: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invite tenantmodels.PupilSectionInvite
		if err := rows.Scan(
			&invite.ID,
			&invite.PupilID,
			&invite.SectionID,
			&invite.InviteDate,
			&invite.Status,
			&invite.PupilFullName,
			&invite.SectionName,
			&invite.PupilEmail,
		); err != nil {
			return nil, fmt.Errorf("error scanning invite row: %v", err)
		}
		invite.TenantID = tenantID
		invite.TenantName = tenantName
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over global invites: %v", err)
	}

	return invites, nil
}

// GetGlobalInvitesForAccount TODO: Add description
func GetGlobalInvitesForAccount(
	accountID int, workspaceDB *sql.DB,
) ([]tenantmodels.GlobalInvite, error) {
	query := `SELECT id, invite_id, account_id, tenant_id
	FROM invite_index WHERE account_id = ?`

	rows, err := workspaceDB.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("error querying global invites for account: %v", err)
	}
	defer rows.Close()

	var invites []tenantmodels.GlobalInvite
	for rows.Next() {
		var invite tenantmodels.GlobalInvite
		if err := rows.Scan(
			&invite.ID,
			&invite.InviteID,
			&invite.AccountID,
			&invite.TenantID,
		); err != nil {
			return nil, fmt.Errorf("error scanning global invite row: %v", err)
		}
		invites = append(invites, invite)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over global invites: %v", err)
	}

	return invites, nil
}

// DeclinePupilSectionInvite TODO: Add description
func DeclinePupilSectionInvite(
	inviteID string, tenantDB *sql.DB,
) error {
	updateQuery := `UPDATE pupils_sections_invite SET status = 'declined' WHERE id = ?`
	_, err := tenantDB.Exec(updateQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error declining pupil section invite: %v", err)
	}

	return nil
}

// AcceptPupilSectionInvite TODO: Add description
func AcceptPupilSectionInvite(
	tenantID int,
	pupilID int,
	inviteID string,
	tenantTx *sql.Tx,
	workspaceTx *sql.Tx,
) error {
	// Update invite status
	updateQuery := `UPDATE pupils_sections_invite SET status = 'accepted' WHERE id = ?`
	_, err := tenantTx.Exec(updateQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error accepting pupil section invite: %v", err)
	}

	// Insert pupil into tenant DB
	insertPupilIntTenantDbQuery := `INSERT IGNORE INTO pupils (id, name, last_name,
	gender, address, guardian_name, phone_number, guardian_number, date_of_birth,
	religion, account_id, place_of_birth)
	SELECT pg.id, pg.name, pg.last_name, pg.gender, pg.address, pg.guardian_name,
	pg.phone_number, pg.guardian_number, pg.date_of_birth, pg.religion,
	pg.account_id, pg.place_of_birth FROM ednevnik_workspace.pupil_global pg
	WHERE pg.id = (SELECT pupil_id FROM pupils_sections_invite WHERE id = ?)`

	_, err = tenantTx.Exec(insertPupilIntTenantDbQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error inserting pupil into tenant DB from invite: %v", err)
	}

	// Insert pupil section relationship
	insertPupilSection := `INSERT INTO pupils_sections (pupil_id, section_id)
	SELECT pupil_id, section_id FROM pupils_sections_invite WHERE id = ?
	ON DUPLICATE KEY UPDATE is_active = 1`

	_, err = tenantTx.Exec(insertPupilSection, inviteID)
	if err != nil {
		return fmt.Errorf("error inserting pupil section from invite: %v", err)
	}

	// Insert into pupil tenant
	insertQuery := `INSERT IGNORE INTO pupil_tenant (pupil_id, tenant_id)
	VALUES (?, ?)`
	_, err = workspaceTx.Exec(insertQuery, pupilID, tenantID)
	if err != nil {
		return fmt.Errorf("error inserting into pupil tenant: %v", err)
	}

	return nil
}

// HandleTeacherSectionAssignmentsHelper TODO: Add description
func HandleTeacherSectionAssignmentsHelper(
	workspaceDB *sql.DB,
	tenantDB *sql.DB,
	teacherID int,
	tenantID int,
	assignmentRequest wpmodels.TeacherSectionAssignment,
) (updatedTeacherAssignmentData []commonmodels.DataForTeacherSectionInvite, err error) {
	// Get teacher account ID
	var teacherAccountID int
	accountQuery := `SELECT account_id FROM teachers WHERE id = ?`
	err = workspaceDB.QueryRow(accountQuery, teacherID).Scan(&teacherAccountID)
	if err != nil {
		return nil, fmt.Errorf("error getting teacher account ID: %v", err)
	}

	// Start a tenantDB transaction
	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// If panic, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	for sectionID, assignment := range assignmentRequest {
		var inviteID int64

		var checkedSubjects []wpmodels.Subject
		for _, subject := range assignment.AvailableSubjects {
			if subject.Checked {
				checkedSubjects = append(checkedSubjects, subject)
			}
		}
		// Check if an invite should be created invite index = 0 (no pending exits)
		// and atleast one pending subjects or teacher pending homeroom is false and
		if assignment.InviteIndexID == 0 && (len(checkedSubjects) > 0 || assignment.HomeroomRequest) {
			// If pending subjects is empty that means an invite does not exist
			// and we need to create a new invite
			// If teacher is a homeroom teacher no need to insert homeroom_teacher value
			var res sql.Result
			if !assignment.IsHomeroom {
				teacherSectionsInviteQuery := `INSERT INTO teachers_sections_invite
				(teacher_id, section_id, homeroom_teacher) VALUES (?, ?, ?)`

				res, err = tx.Exec(teacherSectionsInviteQuery, teacherID, sectionID, assignment.HomeroomRequest)
			} else {
				teacherSectionsInviteQuery := `INSERT INTO teachers_sections_invite
				(teacher_id, section_id) VALUES (?, ?)`

				res, err = tx.Exec(teacherSectionsInviteQuery, teacherID, sectionID)
			}
			if err != nil {
				return nil, fmt.Errorf("error inserting teacher section invite: %v", err)
			}

			inviteID, err = res.LastInsertId()
			if err != nil {
				return nil, fmt.Errorf("error getting last insert ID: %v", err)
			}

			// insert into invite index
			inviteIndexInsert := `INSERT INTO invite_index (invite_id, account_id, tenant_id)
			VALUES (?, ?, ?)`
			_, err = workspaceDB.Exec(inviteIndexInsert, inviteID, teacherAccountID, tenantID)
			if err != nil {
				return nil, fmt.Errorf("error inserting into invite index: %v", err)
			}
		} else {
			// We get the invite ID from the first pending subject
			inviteID = int64(assignment.InviteIndexID)
		}

		// Homeroom handling
		// If teacher is a homeroom teacher and homeroom request if false
		// Delete the assignment
		if !assignment.HomeroomRequest && assignment.IsHomeroom {
			deleteHrQuery := `DELETE FROM homeroom_assignments WHERE section_id = ?`
			_, err := tx.Exec(deleteHrQuery, sectionID)
			if err != nil {
				return nil, fmt.Errorf("error deleting homeroom assignment: %v", err)
			}
		}

		// If homeroom is pending (not assigned yet) and teacher is not already a homeroom
		// teacher just update the status (from NULL which is default) based on homeroom request
		if assignment.InviteIndexID != 0 && !assignment.IsHomeroom {
			pendingHrQuery := `UPDATE teachers_sections_invite SET
			homeroom_teacher = ? WHERE id = ?`
			_, err = tx.Exec(pendingHrQuery, assignment.HomeroomRequest, inviteID)
			if err != nil {
				return nil, fmt.Errorf("error updating pending homeroom: %v", err)
			}
		}

		// Now insert the subjects for this invite
		for _, subject := range checkedSubjects {
			subjectQuery := `INSERT INTO teachers_sections_invite_subjects
				(invite_id, subject_code) VALUES (?, ?)`
			_, err = tx.Exec(subjectQuery, inviteID, subject.SubjectCode)
			if err != nil {
				return nil, fmt.Errorf("error inserting subject for teacher section invite: %v", err)
			}
		}

		// Now for pending subjects that are not checked remove them from the invite subjects
		for _, subject := range assignment.PendingSubjects {
			if !subject.Checked {
				pendingSubjectDeleteQuery := `DELETE FROM teachers_sections_invite_subjects
				WHERE invite_id = ? AND subject_code = ?`
				_, err = tx.Exec(pendingSubjectDeleteQuery, inviteID, subject.SubjectCode)
				if err != nil {
					return nil, fmt.Errorf("error deleting subject for teacher section invite: %v", err)
				}
			}
		}

		// For the assigned subjects that are not checked we need to remove them from the teacher subjects
		for _, subject := range assignment.AssignedSubjects {
			if !subject.Checked {
				assignedSubjectDeleteQuery := `DELETE FROM teachers_sections_subjects
				WHERE teacher_id = ? AND subject_code = ? AND section_id = ?`
				_, err = tx.Exec(assignedSubjectDeleteQuery, teacherID, subject.SubjectCode, sectionID)
				if err != nil {
					return nil, fmt.Errorf("error deleting subject for teacher section: %v", err)
				}
			}
		}

		// Get the count of teacher_sections_invite_subjects for this invite
		countQuery := `SELECT COUNT(*) FROM teachers_sections_invite_subjects
		WHERE invite_id = ?`
		var subjectCount int
		err = tx.QueryRow(countQuery, inviteID).Scan(&subjectCount)
		if err != nil {
			return nil, fmt.Errorf("error counting subjects for teacher section invite: %v", err)
		}

		// If count is 0, we need to delete the invite
		if subjectCount == 0 {
			inviteIndexDeleteQuery := `DELETE FROM ednevnik_workspace.invite_index
			WHERE invite_id = ? AND account_id = ? AND tenant_id = ?
			AND NOT EXISTS (
				SELECT 1 FROM teachers_sections_invite tsi
				WHERE tsi.id = ? AND tsi.homeroom_teacher = 1
			)`

			result, err := tx.Exec(
				inviteIndexDeleteQuery, inviteID, teacherAccountID, tenantID, inviteID,
			)
			if err != nil {
				return nil, fmt.Errorf("error deleting invite index: %v", err)
			}

			// Check if any rows were affected (meaning the invite was not for a homeroom teacher)
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return nil, fmt.Errorf("error checking rows affected: %v", err)
			}

			// Only delete the invite if we successfully deleted from invite_index
			if rowsAffected > 0 {
				deleteInviteQuery := `DELETE FROM teachers_sections_invite WHERE id = ? 
        		AND homeroom_teacher = 0`
				_, err = tx.Exec(deleteInviteQuery, inviteID)
				if err != nil {
					return nil, fmt.Errorf("error deleting teacher section invite: %v", err)
				}
			}
		}

		// Get the count of teachers_sections_subjects for this teacher and section
		teacherSubjectCountQuery := `SELECT COUNT(*) FROM teachers_sections_subjects
		WHERE teacher_id = ? AND section_id = ?`
		var teacherSubjectCount int
		err = tx.QueryRow(teacherSubjectCountQuery, teacherID, sectionID).Scan(&teacherSubjectCount)
		if err != nil {
			return nil, fmt.Errorf("error counting subjects for teacher section: %v", err)
		}
		// If count is 0, delete the teachers sections record
		if teacherSubjectCount == 0 {
			deleteTeacherSectionQuery := `DELETE FROM teachers_sections 
			WHERE teacher_id = ? AND section_id = ? 
			AND NOT EXISTS (
				SELECT 1 FROM homeroom_assignments 
				WHERE teacher_id = ? AND section_id = ?
			)`
			_, err = tx.Exec(deleteTeacherSectionQuery, teacherID, sectionID, teacherID, sectionID)
			if err != nil {
				return nil, fmt.Errorf("error deleting teacher section: %v", err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	updatedData, err := GetDataForTeacherInviteForSingleTeacher(
		strconv.Itoa(teacherID), tenantID, tenantDB, workspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting updated data for teacher invite: %v", err)
	}

	return updatedData, nil
}

func getTeacherInvitesQuery(
	query string,
	args []interface{},
	tenantDB *sql.DB,
	tenantID int,
) ([]wpmodels.TeacherSectionInviteRecord, error) {
	rows, err := tenantDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying teacher invites: %v", err)
	}
	defer rows.Close()

	var invites []wpmodels.TeacherSectionInviteRecord
	var currentInvite *wpmodels.TeacherSectionInviteRecord
	var lastInviteID int

	for rows.Next() {
		var inviteID int
		var teacherID int
		var teacherFullName string
		var sectionID int
		var sectionName string
		var inviteDate string
		var status string
		var subjectCode string
		var subjectName string
		var tenantName string
		var homeroomTeacher bool
		var teacherEmail string

		err := rows.Scan(
			&inviteID,
			&teacherID,
			&teacherFullName,
			&sectionID,
			&sectionName,
			&inviteDate,
			&status,
			&homeroomTeacher,
			&subjectCode,
			&subjectName,
			&tenantName,
			&teacherEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teacher invite row: %v", err)
		}

		// If this is a new invite or first row
		if inviteID != lastInviteID {
			// Add previous invite to slice if exists
			if currentInvite != nil {
				invites = append(invites, *currentInvite)
			}

			// Create new invite
			currentInvite = &wpmodels.TeacherSectionInviteRecord{
				ID:              inviteID,
				TeacherID:       teacherID,
				TeacherFullName: teacherFullName,
				SectionID:       sectionID,
				SectionName:     sectionName,
				InviteDate:      inviteDate,
				Status:          status,
				HomeroomTeacher: homeroomTeacher,
				TenantID:        tenantID,
				TenantName:      tenantName,
				TeacherEmail:    teacherEmail,
				Subjects:        make([]wpmodels.Subject, 0, 5),
			}
			lastInviteID = inviteID
		}

		// Add subject to current invite
		if subjectCode != "" && subjectName != "" {
			subject := wpmodels.Subject{
				SubjectCode: subjectCode,
				SubjectName: subjectName,
			}
			currentInvite.Subjects = append(currentInvite.Subjects, subject)
		}
	}

	// Don't forget the last invite
	if currentInvite != nil {
		invites = append(invites, *currentInvite)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over teacher invite rows: %v", err)
	}

	return invites, nil
}

// GetTeacherInvitesForTenant Get teacher tenant invites data for a single teacher
func GetTeacherInvitesForTenant(
	tenantID,
	inviteID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) ([]wpmodels.TeacherSectionInviteRecord, error) {
	query := `SELECT
		tsi.id,
		tsi.teacher_id,
		CONCAT_WS(' ', t.name, t.last_name) as teacher_full_name,
		tsi.section_id,
		CONCAT(s.class_code, '-', s.section_code) as section_name,
		tsi.invite_date,
		tsi.status,
		tsi.homeroom_teacher,
		COALESCE(sub.subject_code, '') as subject_code,
		COALESCE(sub.subject_name, '') as subject_name, 
		ten.tenant_name,
		a.email
	FROM teachers_sections_invite tsi
	LEFT JOIN teachers_sections_invite_subjects tsis ON tsi.id = tsis.invite_id
	LEFT JOIN ednevnik_workspace.subjects sub ON tsis.subject_code = sub.subject_code
	JOIN sections s ON s.id = tsi.section_id
	JOIN ednevnik_workspace.teachers t ON t.id = tsi.teacher_id
	JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	JOIN ednevnik_workspace.tenant ten ON ten.id = ?
	WHERE tsi.id = ? 
	ORDER BY tsi.id`

	return getTeacherInvitesQuery(query, []interface{}{tenantID, inviteID}, tenantDB, tenantID)
}

// GetAllTeacherInvitesForTenant Get teacher tenant invites data for all teachers in the tenant
func GetAllTeacherInvitesForTenant(
	tenantID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) ([]wpmodels.TeacherSectionInviteRecord, error) {
	query := `SELECT
		tsi.id,
		tsi.teacher_id,
		CONCAT_WS(' ', t.name, t.last_name) as teacher_full_name,
		tsi.section_id,
		CONCAT(s.class_code, '-', s.section_code) as section_name,
		tsi.invite_date,
		tsi.status,
		tsi.homeroom_teacher,
		COALESCE(sub.subject_code, '') as subject_code,
		COALESCE(sub.subject_name, '') as subject_name,
		ten.tenant_name,
		a.email
	FROM teachers_sections_invite tsi
	LEFT JOIN teachers_sections_invite_subjects tsis ON tsi.id = tsis.invite_id
	LEFT JOIN ednevnik_workspace.subjects sub ON tsis.subject_code = sub.subject_code
	JOIN sections s ON s.id = tsi.section_id
	JOIN ednevnik_workspace.teachers t ON t.id = tsi.teacher_id
	JOIN ednevnik_workspace.accounts a ON a.id = t.account_id
	JOIN ednevnik_workspace.tenant ten ON ten.id = ?
	ORDER BY tsi.id`

	return getTeacherInvitesQuery(query, []interface{}{tenantID}, tenantDB, tenantID)
}

// AcceptTeacherSectionInvite TODO: Add description
func AcceptTeacherSectionInvite(
	tenantID int,
	teacherID int,
	inviteID string,
	tenantTx *sql.Tx,
	workspaceTx *sql.Tx,
) (err error) {
	// Update invite status
	updateQuery := `UPDATE teachers_sections_invite SET status = 'accepted' WHERE id = ?`
	_, err = tenantTx.Exec(updateQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error accepting pupil section invite: %v", err)
	}

	// See if teacher already exists for the given sectio
	insertTeacherSectionsRecord := `INSERT IGNORE INTO teachers_sections (teacher_id, section_id)
	SELECT teacher_id, section_id FROM teachers_sections_invite WHERE id = ?`

	_, err = tenantTx.Exec(insertTeacherSectionsRecord, inviteID)
	if err != nil {
		return fmt.Errorf("error inserting teachers_sections record: %v", err)
	}

	teacherSubjectsInsertQuery := `INSERT INTO teachers_sections_subjects
	(section_id, subject_code, teacher_id)
	SELECT tsi.section_id, tsis.subject_code, tsi.teacher_id
	FROM teachers_sections_invite_subjects tsis
	JOIN teachers_sections_invite tsi ON tsi.id = tsis.invite_id
	WHERE tsi.id = ?`

	_, err = tenantTx.Exec(teacherSubjectsInsertQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error inserting subjects for teacher: %v", err)
	}

	// Handle homeroom teacher assignment
	// First, check if this invite is for a homeroom teacher
	var isHomeroomTeacher bool
	var sectionID int
	homeroomCheckQuery := `SELECT homeroom_teacher, section_id FROM teachers_sections_invite WHERE id = ?`
	err = tenantTx.QueryRow(homeroomCheckQuery, inviteID).Scan(&isHomeroomTeacher, &sectionID)
	if err != nil {
		return fmt.Errorf("error checking homeroom teacher status: %v", err)
	}

	// If this teacher is being assigned as homeroom teacher
	if isHomeroomTeacher {
		// Use INSERT ... ON DUPLICATE KEY UPDATE to either insert new or update existing
		homeroomAssignmentQuery := `INSERT INTO homeroom_assignments (section_id, teacher_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE teacher_id = VALUES(teacher_id)`

		_, err = tenantTx.Exec(homeroomAssignmentQuery, sectionID, teacherID)
		if err != nil {
			return fmt.Errorf("error inserting/updating homeroom assignment: %v", err)
		}
	}

	teacherTenantInsertQuery := `INSERT IGNORE INTO teacher_tenant (teacher_id, tenant_id)
	VALUES (?, ?)`
	_, err = workspaceTx.Exec(teacherTenantInsertQuery, teacherID, tenantID)
	if err != nil {
		return fmt.Errorf("error inserting teacher tenant record: %v", err)
	}

	return nil
}

// DeclineTeacherSectionInvite TODO: Add description
func DeclineTeacherSectionInvite(
	inviteID string,
	tenantDB *sql.DB,
) (err error) {
	// Update invite status
	updateQuery := `UPDATE teachers_sections_invite SET status = 'declined' WHERE id = ?`
	_, err = tenantDB.Exec(updateQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error accepting teacher section invite: %v", err)
	}

	return nil
}

// DeletePupilInviteHelper TODO: Add description
func DeletePupilInviteHelper(
	inviteID,
	pupilID,
	tenantID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) error {
	pupilSectionInviteDeleteQuery := `DELETE FROM pupils_sections_invite
	WHERE id = ?`
	_, err := tenantDB.Exec(pupilSectionInviteDeleteQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error deleting pupil_sections_invite record: %v", err)
	}

	var accountID int
	getPupilAccountQuery := `SELECT account_id FROM pupil_global WHERE
	id = ?`
	err = workspaceDB.QueryRow(getPupilAccountQuery, pupilID).Scan(
		&accountID,
	)
	if err != nil {
		return fmt.Errorf("error getting pupil account id: %v", err)
	}

	globalInviteDeleteQuery := `DELETE FROM invite_index WHERE
	invite_id = ? AND account_id = ? AND tenant_id = ?`
	_, err = workspaceDB.Exec(
		globalInviteDeleteQuery, inviteID, accountID, tenantID,
	)
	if err != nil {
		return fmt.Errorf("error deleting global pupil invite: %v", err)
	}

	return nil
}

// DeleteTeacherInviteHelper TODO: Add description
func DeleteTeacherInviteHelper(
	inviteID,
	teacherID,
	tenantID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) error {
	teacherSectionInviteDeleteQuery := `DELETE FROM teachers_sections_invite
	WHERE id = ?`
	_, err := tenantDB.Exec(teacherSectionInviteDeleteQuery, inviteID)
	if err != nil {
		return fmt.Errorf("error deleting teacher_sections_invite record: %v", err)
	}

	var accountID int
	getTeacherAccountQuery := `SELECT account_id FROM teachers WHERE
	id = ?`
	err = workspaceDB.QueryRow(getTeacherAccountQuery, teacherID).Scan(
		&accountID,
	)
	if err != nil {
		return fmt.Errorf("error getting teacher account id: %v", err)
	}

	globalInviteDeleteQuery := `DELETE FROM invite_index WHERE
	invite_id = ? AND account_id = ? AND tenant_id = ?`
	_, err = workspaceDB.Exec(
		globalInviteDeleteQuery, inviteID, accountID, tenantID,
	)
	if err != nil {
		return fmt.Errorf("error deleting global teacher invite: %v", err)
	}

	return nil
}
