package tenantfactory

import (
	tenantmodels "ednevnik-backend/models/tenant"
	"ednevnik-backend/util"
)

// CreateSchedule TODO: Add description
func (t *ConfigurableTenant) CreateSchedule(
	data tenantmodels.ScheduleGroupCollection, sectionID string,
) error {
	err := util.CreateSchedule(
		data, t.UserTenantDB, sectionID,
	)

	return err
}

// GetScheduleForSection TODO: Add description
func (t *ConfigurableTenant) GetScheduleForSection(
	sectionID string,
) (tenantmodels.ScheduleGroupCollection, error) {
	schedule, err := util.GetScheduleForSection(
		sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}

// GetScheduleForTeacher TODO: Add description
func (t *ConfigurableTenant) GetScheduleForTeacher(
	teacherID string,
) (tenantmodels.ScheduleGroupCollection, error) {
	schedule, err := util.GetScheduleForTeacher(
		teacherID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	return schedule, nil
}
