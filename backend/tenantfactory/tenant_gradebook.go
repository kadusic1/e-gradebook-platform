package tenantfactory

import (
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	"ednevnik-backend/util"
	"fmt"
)

// GetSectionGradesForSubject retrieves grades for all pupils in a section for a specific subject.
func (t *ConfigurableTenant) GetSectionGradesForSubject(
	sectionID int,
	semesterCode,
	subjectCode string,
) ([]tenantmodels.GradePupilGroup, error) {
	grades, err := util.GetPupilGradesForSectionSubject(
		sectionID,
		semesterCode,
		subjectCode,
		t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	return grades, nil
}

// CreateGrade inserts a new grade into the database and returns the created grade.
func (t *ConfigurableTenant) CreateGrade(
	grade *tenantmodels.Grade,
) (*tenantmodels.GradePupilGroup, error) {
	createdGrade, err := util.CreateGrade(
		grade,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}

	return createdGrade, nil
}

// DeleteGrade removes a grade from the database and returns the updated grades for the pupil.
func (t *ConfigurableTenant) DeleteGrade(
	grade *tenantmodels.Grade, teacherID int,
) (*tenantmodels.GradePupilGroup, error) {
	teacherForSignature, err := util.GetTeacherByID(
		fmt.Sprintf("%d", teacherID), t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}
	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	gradesAfterDeletion, err := util.DeleteGrade(
		grade,
		signature,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}
	return gradesAfterDeletion, err
}

// UpdateGrade updates an existing grade
func (t *ConfigurableTenant) UpdateGrade(
	grade *tenantmodels.Grade,
) (*tenantmodels.GradePupilGroup, error) {
	createdGrade, err := util.EditGrade(
		grade,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}

	return createdGrade, nil
}

// GetPupilGradesForSectionPupil returns subject grades for a pupil in a section
func (t *ConfigurableTenant) GetPupilGradesForSectionPupil(
	sectionID,
	pupilID int,
	semesterCode string,
) ([]commonmodels.GradeSubjectGroup, error) {
	grades, err := util.GetPupilGradesForSectionPupilHelper(
		sectionID,
		pupilID,
		semesterCode,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}
	return grades, nil
}

// GetSectionBehaviourGradesForPupil uses the GetSectionBehaviourGradesForPupilHelper
// util function to return pupil behaviour grades for a specific section.
func (t *ConfigurableTenant) GetSectionBehaviourGradesForPupil(
	pupilID, sectionID int,
) ([]tenantmodels.BehaviourGrade, error) {
	behaviourGrades, err := util.GetSectionBehaviourGradesForPupilHelper(
		pupilID, sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	return behaviourGrades, nil
}

// GetGradeEditHistory retrieves the edit history of a specific grade by its ID.
func (t *ConfigurableTenant) GetGradeEditHistory(gradeID int) ([]tenantmodels.Grade, error) {
	grades, err := util.GetGradeEditHistoryHelper(
		gradeID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	return grades, nil
}

// GetBehaviourGradeHistory retrieves the history of a specific behaviour grade by its ID.
func (t *ConfigurableTenant) GetBehaviourGradeHistory(
	behaviourGradeID int,
) ([]tenantmodels.BehaviourGrade, error) {
	behaviourGrades, err := util.GetPupilBehaviourGradeHistoryHelper(
		behaviourGradeID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	return behaviourGrades, nil
}

// GetCompleteGradebookData retrieves all gradebook data for a section.
func (t *ConfigurableTenant) GetCompleteGradebookData(
	sectionID int,
) (*tenantmodels.CompleteGradebook, error) {
	var completeGradebook tenantmodels.CompleteGradebook

	pupils, err := util.GetPupilsForCompleteGradebook(
		fmt.Sprintf("%d", sectionID), t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.Pupils = pupils

	scheduleHistory, err := util.GetAllSchedulesForSection(
		sectionID, t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.ScheduleHistory = scheduleHistory

	sectionSemesters, err := t.GetSemestersForSection(
		fmt.Sprintf("%d", sectionID),
	)
	if err != nil {
		return nil, err
	}

	section, err := util.GetSectionByID(
		int64(sectionID), t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}

	sectionSubjects, err := util.GetAllSubjectsForCurriculumCode(
		section.CurriculumCode, t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.Subjects = sectionSubjects

	gradeData, err := util.GetGradeDataForCompleteGradebook(
		sectionID,
		sectionSemesters,
		pupils,
		sectionSubjects,
		t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.GradeData = gradeData

	lessons, err := util.GetLessonsByWeekForSection(
		sectionID,
		t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.Lessons = lessons

	absences, err := util.GetAbsencesByWeekForSection(
		sectionID,
		t.UserTenantDB,
	)
	if err != nil {
		return nil, err
	}
	completeGradebook.Absences = absences

	return &completeGradebook, nil
}
