package tenantmodels

import (
	wpmodels "ednevnik-backend/models/workspace"
)

// SectionCreateMetadata TODO: Add description
type SectionCreateMetadata struct {
	Classes     []wpmodels.Class         `json:"classes"`
	Teachers    []wpmodels.Teacher       `json:"teachers"`
	Curriculums []wpmodels.CurriculumGet `json:"curriculums"`
}

// SectionCreate TODO: Add description
type SectionCreate struct {
	SectionCode    string `json:"section_code"`
	ClassCode      string `json:"class_code"`
	Year           string `json:"year"`
	CurriculumCode string `json:"curriculum_code"`
}

// Section TODO: Add description
type Section struct {
	// Fields from the database
	ID             int64  `json:"id,omitempty"`
	SectionCode    string `json:"section_code"`
	ClassCode      string `json:"class_code"`
	Year           string `json:"year"`
	TenantID       int    `json:"tenant_id,omitempty"`
	CurriculumCode string `json:"curriculum_code"`
	Archived       bool   `json:"archived,omitempty"`
	// Additional fields for response
	HomeroomTeacherID       int    `json:"homeroom_teacher_id,omitempty"`
	Name                    string `json:"name,omitempty"`
	HomeroomTeacherEmail    string `json:"homeroom_teacher_email,omitempty"`
	HomeroomTeacherFullName string `json:"homeroom_teacher_full_name,omitempty"`
	CurriculumName          string `json:"curriculum_name,omitempty"`
	// Addition color config field
	ColorConfig string `json:"color_config,omitempty"`
	// Optional tenant_name field when returning sections for pupils
	TenantName         string `json:"tenant_name,omitempty"`
	PupilDisplay       string `json:"pupil_display,omitempty"`
	PupilInviteDisplay string `json:"pupil_invite_display,omitempty"`
	LessonDisplay      string `json:"lesson_display,omitempty"`
	AbsenceDisplay     string `json:"absence_display,omitempty"`
	// Just for secondary schools
	CourseName string `json:"course_name,omitempty"`
}
