package tenantfactory

import (
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	"ednevnik-backend/util"
	"fmt"
)

// GetPupilsForSection TODO: Add description
func (t *ConfigurableTenant) GetPupilsForSection(
	sectionID string, includeUnenrolled bool,
) (*commonmodels.GetSectionPupilsResponse, error) {

	pupils, err := util.GetPupilsForSection(
		sectionID, includeUnenrolled, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	pendingPupils, err := util.GetPupilSectionInvitesForSectionHelper(
		sectionID, t.TenantData.TenantName, t.UserTenantDB, int(t.TenantData.ID),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting pending pupils for section: %v", err,
		)
	}

	pupilsForAssignment, err := util.GetPupilsForTenantSectionAssignment(
		sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error getting pupils for tenant section assignment: %v", err,
		)
	}

	response := commonmodels.GetSectionPupilsResponse{
		Pupils:              pupils,
		PendingPupils:       pendingPupils,
		PupilsForAssignment: pupilsForAssignment,
	}

	return &response, nil
}

// DeletePupilFromSection TODO: Add description
func (t *ConfigurableTenant) DeletePupilFromSection(pupilID, sectionID string) error {
	// First check how many pupil tenant records are there
	var pupilSectionCount int

	pupilSectionCountQuery := `SELECT COUNT(*) FROM pupils_sections ps
		JOIN sections s ON ps.section_id = s.id
		WHERE ps.pupil_id = ? AND s.tenant_id = ?`

	err := t.UserTenantDB.QueryRow(pupilSectionCountQuery, pupilID, t.TenantData.ID).Scan(
		&pupilSectionCount,
	)
	if err != nil {
		return fmt.Errorf("error getting pupil section count: %v", err)
	}

	if pupilSectionCount < 1 {
		return fmt.Errorf("pupil does not exist in any sections for this tenant")
	}

	if pupilSectionCount == 1 {
		err = util.DeleteTenantPupilRecord(
			pupilID, int(t.TenantData.ID), t.UserWorkspaceDB, t.UserTenantDB,
		)
		if err != nil {
			return fmt.Errorf("error deleting pupil tenat record: %v", err)
		}
	} else {
		pupilSectionDeleteQuery := `DELETE FROM pupils_sections
			WHERE pupil_id = ? AND section_id = ?`
		_, err := t.UserTenantDB.Exec(
			pupilSectionDeleteQuery, pupilID, sectionID,
		)
		if err != nil {
			return fmt.Errorf("error deleting pupil from section: %v", err)
		}
	}

	return nil
}

// UpdatePupil TODO: Add description
func (t *ConfigurableTenant) UpdatePupil(
	oldPupil tenantmodels.Pupil, newPupil tenantmodels.Pupil,
) error {

	err := util.UpdatePupilTenantRecord(
		fmt.Sprintf("%d", oldPupil.ID),
		newPupil,
		t.UserTenantDB,
	)
	if err != nil {
		return fmt.Errorf(
			"error updating pupil for tenant: %v", err,
		)
	}

	return nil
}

// DeletePupilFromTenant TODO: Add description
func (t *ConfigurableTenant) DeletePupilFromTenant(pupilID string) error {
	err := util.DeleteTenantPupilRecord(
		pupilID, int(t.TenantData.ID), t.UserWorkspaceDB, t.UserTenantDB,
	)

	if err != nil {
		return fmt.Errorf("error deleting pupil from tenant: %v", err)
	}

	return nil
}

// UpdatePupilBehaviourGrade updates a pupil behaviour grade for a pupil in a
// section
func (t *ConfigurableTenant) UpdatePupilBehaviourGrade(
	behaviourGradesToUpdate tenantmodels.BehaviourGrade, teacherID int,
) (*tenantmodels.BehaviourGrade, error) {
	teacherForSignature, err := util.GetTeacherByID(
		fmt.Sprintf("%d", teacherID), t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting teacher for signature: %v", err)
	}

	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	behaviourGrade, err := util.UpdatePupilBehaviourGradeHelper(
		signature,
		behaviourGradesToUpdate,
		t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	return behaviourGrade, nil
}
