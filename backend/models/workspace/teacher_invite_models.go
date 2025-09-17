package wpmodels

// TeacherSectionAssignmentRecord TODO: Add description
type TeacherSectionAssignmentRecord struct {
	AvailableSubjects []Subject `json:"availableSubjects"`
	PendingSubjects   []Subject `json:"pendingSubjects"`
	AssignedSubjects  []Subject `json:"assignedSubjects"`
	IsHomeroom        bool      `json:"isHomeroom"`
	PendingHomeroom   bool      `json:"pendingHomeroom"`
	InviteIndexID     int       `json:"inviteIndexID"`
	HomeroomRequest   bool      `json:"homeroomRequest"`
}

// TeacherSectionAssignment is a map of section IDs to their respective TeacherSectionInviteRecord
type TeacherSectionAssignment map[string]TeacherSectionAssignmentRecord

// TeacherSectionInviteRecord TODO: Add description
type TeacherSectionInviteRecord struct {
	ID              int       `json:"id"`
	TeacherID       int       `json:"teacher_id"`
	TeacherFullName string    `json:"teacher_full_name"`
	SectionID       int       `json:"section_id"`
	SectionName     string    `json:"section_name"`
	InviteDate      string    `json:"invite_date"`
	Status          string    `json:"status"`
	Subjects        []Subject `json:"subjects"`
	TenantID        int       `json:"tenant_id,omitempty"`
	TenantName      string    `json:"tenant_name,omitempty"`
	HomeroomTeacher bool      `json:"homeroom_teacher"`
	TeacherEmail    string    `json:"teacher_email,omitempty"`
}
