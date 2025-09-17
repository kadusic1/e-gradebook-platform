package tenantfactory

import (
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"
	"fmt"
)

// GetLessonsForSection retrieves all lessons associated with a specific section ID
// from the tenant's database. Returns a slice of LessonData or an error if the
// operation fails.
func (t *ConfigurableTenant) GetLessonsForSection(
	sectionID int, claims *wpmodels.Claims,
) ([]tenantmodels.LessonData, error) {
	lessons, err := util.GetLessonsForSection(
		sectionID, t.UserTenantDB, claims,
	)
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

// CreateSectionLesson creates a new lesson in the tenant's database using the provided
// lesson data. Returns an error if the creation operation fails.
func (t *ConfigurableTenant) CreateSectionLesson(
	requestData tenantmodels.LessonData, teacherID int,
) (*tenantmodels.LessonData, error) {
	teacherForSignature, err := util.GetTeacherByID(
		fmt.Sprintf("%d", teacherID),
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}

	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	newLesson, err := util.CreateLesson(
		requestData, t.UserTenantDB, signature,
	)
	if err != nil {
		return nil, err
	}
	return newLesson, nil
}

// UpdateLesson updates an existing lesson identified by lessonID in the tenant's database
// with the provided lesson data. Returns an error if the update operation fails.
func (t *ConfigurableTenant) UpdateLesson(
	lessonID int, requestData tenantmodels.LessonData, teacherID int,
) (*tenantmodels.LessonData, error) {
	teacherForSignature, err := util.GetTeacherByID(
		fmt.Sprintf("%d", teacherID),
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}

	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	updatedLesson, err := util.UpdateLesson(
		lessonID, requestData, t.UserTenantDB, signature,
	)
	if err != nil {
		return nil, err
	}
	return updatedLesson, err
}

// DeleteLesson removes a lesson identified by lessonID from the tenant's database.
// Returns an error if the deletion operation fails.
func (t *ConfigurableTenant) DeleteLesson(
	lessonID int,
) error {
	err := util.DeleteLesson(
		lessonID, t.UserTenantDB,
	)
	return err
}

// GetLessonByID retrieves a lesson by its ID from the tenant's database.
func (t *ConfigurableTenant) GetLessonByID(
	lessonID int,
) (*tenantmodels.LessonData, error) {
	lesson, err := util.GetLessonByID(lessonID, t.UserTenantDB)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

// GetAbsentAttendancesForSection retrieves all absent, excused, and unexcused
// pupil attendance records for a specific section ID from the tenant's database.
func (t *ConfigurableTenant) GetAbsentAttendancesForSection(
	sectionID int,
) ([]tenantmodels.PupilAttendance, error) {
	attendances, err := util.GetAbsentAttendancesForSectionHelper(
		sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// GetAbsentAttendancesForPupil retrieves all absent, excused, and unexcused
// pupil attendance records for a specific pupil ID from the tenant's database.
func (t *ConfigurableTenant) GetAbsentAttendancesForPupil(
	pupilID, sectionID int,
) ([]tenantmodels.PupilAttendance, error) {
	attendances, err := util.GetAbsentAttendancesForPupilHelper(
		pupilID, sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// HandleAttendanceAction processes attendance actions such as marking a pupil's
// absence as excused or unexcused.
func (t *ConfigurableTenant) HandleAttendanceAction(
	action tenantmodels.AttendanceAction,
) error {
	err := util.HandleAttendanceActionHelper(
		action, t.UserTenantDB,
	)
	return err
}
