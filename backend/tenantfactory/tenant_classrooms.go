package tenantfactory

import (
	tenantmodels "ednevnik-backend/models/tenant"
	"ednevnik-backend/util"
)

// CreateClassroom TODO: Add description
func (t *ConfigurableTenant) CreateClassroom(
	data tenantmodels.Classroom,
) error {
	return util.CreateClassroom(data, t.UserTenantDB)
}

// UpdateClassroom TODO: Add description
func (t *ConfigurableTenant) UpdateClassroom(
	data tenantmodels.Classroom, oldCode string,
) error {
	return util.UpdateClassroom(data, t.UserTenantDB, oldCode)
}

// GetAllClassroomsForTenant TODO: Add description
func (t *ConfigurableTenant) GetAllClassroomsForTenant() ([]tenantmodels.Classroom, error) {
	return util.GetAllClassroomsForTenant(t.UserTenantDB)
}

// DeleteClassroom TODO: Add description
func (t *ConfigurableTenant) DeleteClassroom(code string) error {
	return util.DeleteClassroom(code, t.UserTenantDB)
}
