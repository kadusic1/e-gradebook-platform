package tenantfactory

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
)

// SendPupilSectionInvite TODO: Add description
func (t *ConfigurableTenant) SendPupilSectionInvite(
	pupilID int, sectionID, tenantID string,
) (*tenantmodels.PupilSectionInvite, error) {

	newInvite, err := util.SendPupilSectionInvite(
		pupilID, sectionID, tenantID, t.UserTenantDB, t.UserWorkspaceDB, t.TenantData.TenantName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error sending pupil section invite: %v", err,
		)
	}

	return newInvite, nil
}

// GetPupilSectionInvite TODO: Add description
func (t *ConfigurableTenant) GetPupilSectionInvite(
	inviteID int,
) (*tenantmodels.PupilSectionInvite, error) {

	pupilSectionInvites, err := util.GetPupilSectionInvite(
		inviteID, t.TenantData.TenantName, t.UserTenantDB, int(t.TenantData.ID),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting pupil section invites: %v", err,
		)
	}

	return pupilSectionInvites, nil
}

// AcceptPupilSectionInvite TODO: Add description
func (t *ConfigurableTenant) AcceptPupilSectionInvite(
	inviteID string,
) error {

	var err error
	// Start transactions
	workspaceTx, tenantTx, err := t.StartTransactions(
		t.UserWorkspaceDB, t.UserTenantDB,
	)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tenantTx.Rollback()
			_ = workspaceTx.Rollback()
		}
	}()

	// Read pupil id from invite
	var pupilID int
	pupilQuery := `SELECT pupil_id FROM pupils_sections_invite WHERE id = ?`
	err = tenantTx.QueryRow(pupilQuery, inviteID).Scan(&pupilID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invite not found for ID: %s", inviteID)
		}
		return fmt.Errorf("error retrieving pupil ID from invite: %v", err)
	}

	err = util.AcceptPupilSectionInvite(
		int(t.TenantData.ID), pupilID, inviteID, tenantTx, workspaceTx,
	)
	if err != nil {
		return fmt.Errorf(
			"error accepting pupil section invite: %v", err,
		)
	}

	if err = tenantTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tenant DB transaction: %w", err)
	}

	if err = workspaceTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit workspace DB transaction: %w", err)
	}

	return nil
}

// AcceptTeacherSectionInvite TODO: Add description
func (t *ConfigurableTenant) AcceptTeacherSectionInvite(
	inviteID string,
) error {
	var err error
	// Start transactions
	workspaceTx, tenantTx, err := t.StartTransactions(
		t.UserWorkspaceDB, t.UserTenantDB,
	)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tenantTx.Rollback()
			_ = workspaceTx.Rollback()
		}
	}()

	// Read teacher id from invite
	var teacherID int
	teacherQuery := `SELECT teacher_id FROM teachers_sections_invite WHERE id = ?`
	err = tenantTx.QueryRow(teacherQuery, inviteID).Scan(&teacherID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invite not found for ID: %s", inviteID)
		}
		return fmt.Errorf("error retrieving teacher ID from invite: %v", err)
	}

	err = util.AcceptTeacherSectionInvite(
		int(t.TenantData.ID), teacherID, inviteID, tenantTx, workspaceTx,
	)
	if err != nil {
		return fmt.Errorf(
			"error accepting teacher section invite: %v", err,
		)
	}

	if err = tenantTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tenant DB transaction: %w", err)
	}

	if err = workspaceTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit workspace DB transaction: %w", err)
	}

	return nil
}

// DeclineTeacherSectionInvite TODO: Add description
func (t *ConfigurableTenant) DeclineTeacherSectionInvite(
	inviteID string,
) error {
	err := util.DeclineTeacherSectionInvite(inviteID, t.UserTenantDB)
	if err != nil {
		return fmt.Errorf("error declining teacher invite: %v", err)
	}

	return nil
}

// DeclinePupilSectionInvite TODO: Add description
func (t *ConfigurableTenant) DeclinePupilSectionInvite(inviteID string) error {

	err := util.DeclinePupilSectionInvite(
		inviteID, t.UserTenantDB,
	)
	if err != nil {
		return fmt.Errorf(
			"error declining pupil section invite: %v", err,
		)
	}

	return nil
}

// GetDataForTeacherInviteForTenant TODO: Add description
func (t *ConfigurableTenant) GetDataForTeacherInviteForTenant() ([]commonmodels.DataForTeacherSectionInvite, error) {
	data, err := util.GetDataForTeacherInviteForTenantHelper(
		int(t.TenantData.ID), t.UserTenantDB, t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting data for teacher invite for tenant: %v", err,
		)
	}

	return data, nil
}

// HandleTeacherSectionAssignments TODO: Add description
func (t *ConfigurableTenant) HandleTeacherSectionAssignments(
	teacherID int, assignmentRequest wpmodels.TeacherSectionAssignment,
) ([]commonmodels.DataForTeacherSectionInvite, []wpmodels.Teacher, error) {
	data, err := util.HandleTeacherSectionAssignmentsHelper(
		t.UserWorkspaceDB, t.UserTenantDB, teacherID, int(t.TenantData.ID), assignmentRequest,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"error handling teacher section assignments: %v", err,
		)
	}

	// Get the count of teacher section records
	var teacherSectionCount int
	teacherSectionsCountQuery := `SELECT COUNT(*) FROM teachers_sections
	WHERE teacher_id = ?`
	err = t.UserTenantDB.QueryRow(teacherSectionsCountQuery, teacherID).Scan(&teacherSectionCount)
	if err != nil {
		return nil, nil, fmt.Errorf("error counting sections for teacher: %v", err)
	}

	var teachers []wpmodels.Teacher

	// If count is 0, delete the teacher tenant record
	if teacherSectionCount == 0 {
		deleteTeacherTenantQuery := `DELETE FROM teacher_tenant WHERE
		teacher_id = ? AND tenant_id = ?`
		res, err := t.UserWorkspaceDB.Exec(deleteTeacherTenantQuery, teacherID, t.TenantData.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("error deleting teacher tenant record: %v", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, nil, fmt.Errorf("error checking rows affected: %v", err)
		}

		if rowsAffected > 0 {
			// Get teachers for tenant
			teachers, err = t.GetTeachersForTenant()
			if err != nil {
				return nil, nil, fmt.Errorf("error getting teachers for tenant: %v", err)
			}
		}
	}

	return data, teachers, nil
}

// GetInvitesForTeacher TODO: Add description
func (t *ConfigurableTenant) GetInvitesForTeacher(
	inviteID int,
) ([]wpmodels.TeacherSectionInviteRecord, error) {
	invites, err := util.GetTeacherInvitesForTenant(
		int(t.TenantData.ID), inviteID, t.UserTenantDB, t.UserWorkspaceDB,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting tenant invites for teacher: %v", err)
	}

	return invites, nil
}

// GetAllTeacherInvites TODO: Add description
func (t *ConfigurableTenant) GetAllTeacherInvites() (
	[]wpmodels.TeacherSectionInviteRecord, error,
) {
	invites, err := util.GetAllTeacherInvitesForTenant(
		int(t.TenantData.ID), t.UserTenantDB, t.UserWorkspaceDB,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting tenant invites: %v", err)
	}

	return invites, nil
}

// DeletePupilInvite TODO: Add description
func (t *ConfigurableTenant) DeletePupilInvite(inviteID, pupilID int) error {
	err := util.DeletePupilInviteHelper(
		inviteID,
		pupilID,
		int(t.TenantData.ID),
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)

	if err != nil {
		return fmt.Errorf("error deleting pupil invite: %v", err)
	}

	return nil
}

// DeleteTeacherInvite TODO: Add description
func (t *ConfigurableTenant) DeleteTeacherInvite(inviteID, teacherID int) error {
	err := util.DeleteTeacherInviteHelper(
		inviteID,
		teacherID,
		int(t.TenantData.ID),
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)

	if err != nil {
		return fmt.Errorf("error deleting teacher invite: %v", err)
	}

	return nil
}
