package commonmodels

import (
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
)

// GetSectionPupilsResponse TODO: Add description
type GetSectionPupilsResponse struct {
	Pupils              []tenantmodels.Pupil              `json:"pupils"`
	PendingPupils       []tenantmodels.PupilSectionInvite `json:"pending_pupils"`
	PupilsForAssignment []tenantmodels.Pupil              `json:"pupils_for_assignment"`
}

// DataForTeacherSectionInvite TODO: Add description
type DataForTeacherSectionInvite struct {
	Section                tenantmodels.Section `json:"section"`
	Teacher                wpmodels.Teacher     `json:"teacher"`
	AllSubjects            []wpmodels.Subject   `json:"all_subjects"`
	AssignedSubjects       []wpmodels.Subject   `json:"assigned_subjects"`
	PendingSubjects        []wpmodels.Subject   `json:"pending_subjects"`
	AvailableSubjects      []wpmodels.Subject   `json:"available_subjects"`
	IsHomeroomTeacher      bool                 `json:"is_homeroom_teacher"`
	PendingHomeroomTeacher bool                 `json:"is_pending_homeroom_teacher"`
	InviteIndexID          int                  `json:"invite_index_id"`
}

// GradeSubjectGroup represents a group of grades for a subject
// Used to display results for pupils
type GradeSubjectGroup struct {
	Subject      wpmodels.Subject     `json:"subject"`
	AverageGrade float64              `json:"average_grade,omitempty"`
	Grades       []tenantmodels.Grade `json:"grades"`
}

// PupilSubjectSemester represents a pupil and subject combination
// Eg. used to return subjects without finalized grades for a pupil (1 to 1
// mapping)
type PupilSubjectSemester struct {
	Subject  wpmodels.Subject        `json:"subject"`
	Pupil    tenantmodels.Pupil      `json:"pupil"`
	Semester wpmodels.TenantSemester `json:"semester"`
}

// Certificate represents a document that contains pupil grades, behaviour
// grades, and other relevant information.
type Certificate struct {
	Tenant         wpmodels.Tenant             `json:"tenant"`
	Section        tenantmodels.Section        `json:"section"`
	Pupil          tenantmodels.Pupil          `json:"pupil"`
	FinalGrades    []tenantmodels.Grade        `json:"final_grades"`
	BehaviourGrade tenantmodels.BehaviourGrade `json:"behaviour_grades"`
	AverageGrade   float64                     `json:"average_grade,omitempty"`
	GraduateGrade  int                         `json:"graduate_grade,omitempty"`
	Passed         bool                        `json:"passed"`
	// Just for secondary schools
	CourseName string `json:"course_name,omitempty"`
}

type AIPermissionData struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	LastName            string   `json:"last_name"`
	Email               string   `json:"email"`
	Phone               string   `json:"phone"`
	AccountType         string   `json:"account_type"`
	AccountID           int      `json:"account_id"`
	TenantIDs           []string `json:"tenant_ids"`
	TenantAdminTenantID int      `json:"tenant_admin_tenant_id,omitempty"`
	TenantNames         []string `json:"tenant_names,omitempty"`
}

// ChatRequest represents the request payload for chatbot
type ChatRequest struct {
	Question       string           `json:"question"`
	SessionID      string           `json:"session_id,omitempty"`
	PermissionData AIPermissionData `json:"permission_data,omitempty"`
}

// ChatResponse represents the response from chatbot
type ChatResponse struct {
	Answer    string `json:"answer"`
	SessionID string `json:"session_id"`
}
