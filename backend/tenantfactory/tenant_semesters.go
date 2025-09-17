package tenantfactory

import (
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
)

// UpdateTenantSemesterDates TODO: Add description
func (t *ConfigurableTenant) UpdateTenantSemesterDates(
	semesterCode, startDate, endDate, nppCode string,
) (wpmodels.TenantSemester, error) {
	updatedSemester, err := util.UpdateTenantSemesterDates(
		t.UserWorkspaceDB,
		fmt.Sprintf("%d", t.TenantData.ID),
		semesterCode, startDate, endDate, nppCode,
	)
	if err != nil {
		return wpmodels.TenantSemester{}, fmt.Errorf(
			"failed to update semester %s: %w", semesterCode, err,
		)
	}
	return updatedSemester, nil
}

// GetSemestersForTenant TODO: Add description
func (t *ConfigurableTenant) GetSemestersForTenant() ([]wpmodels.TenantSemester, error) {
	semesters, err := util.GetSemestersForTenant(
		t.UserWorkspaceDB,
		fmt.Sprintf("%d", t.TenantData.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting semesters for tenant: %v", err)
	}

	return semesters, nil
}

// GetSemestersForSection returns semesters (polugodišta) for a section
func (t *ConfigurableTenant) GetSemestersForSection(
	sectionID string,
) ([]wpmodels.TenantSemester, error) {
	sectionSemesters, err := util.GetSemestersForSectionHelper(
		t.UserWorkspaceDB,
		t.UserTenantDB,
		fmt.Sprintf("%d", t.TenantData.ID),
		sectionID,
	)
	if err != nil {
		return nil, err
	}
	return sectionSemesters, nil
}
