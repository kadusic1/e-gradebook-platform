package wpmodels

// Tenant represents a tenant in the system
type Tenant struct {
	ID                   int64    `json:"id"`
	TenantName           string   `json:"tenant_name"`
	CantonCode           string   `json:"canton_code"`
	Address              string   `json:"address"`
	Phone                string   `json:"phone"`
	Email                string   `json:"email"`
	DirectorName         string   `json:"director_name"`
	TenantType           string   `json:"tenant_type"`
	Longitude            *float64 `json:"longitude,omitempty"`
	Latitude             *float64 `json:"latitude,omitempty"`
	AIEnabled            *bool    `json:"ai_enabled,omitempty"`
	Domain               *string  `json:"domain,omitempty"`
	TeacherName          string   `json:"teacher_name,omitempty"`
	TeacherLastName      string   `json:"teacher_last_name,omitempty"`
	TeacherEmail         string   `json:"teacher_email,omitempty"`
	TeacherPhone         string   `json:"teacher_phone,omitempty"`
	TeacherPassword      string   `json:"teacher_password,omitempty"`
	TeacherID            int      `json:"teacher_id,omitempty"`
	ColorConfig          string   `json:"color_config"`
	TeacherDisplay       string   `json:"teacher_display,omitempty"`
	TeacherInviteDisplay string   `json:"teacher_invite_display,omitempty"`
	PupilDisplay         string   `json:"pupil_display,omitempty"`
	PupilInviteDisplay   string   `json:"pupil_invite_display,omitempty"`
	SectionDisplay       string   `json:"section_display,omitempty"`
	CurriculumDisplay    string   `json:"curriculum_display,omitempty"`
	SemesterDisplay      string   `json:"semester_display,omitempty"`
	LessonDisplay        string   `json:"lesson_display,omitempty"`
	AbsenceDisplay       string   `json:"absence_display,omitempty"`
	ClassroomDisplay     string   `json:"classroom_display,omitempty"`
	CantonName           string   `json:"canton_name,omitempty"`
	TenantCity           string   `json:"tenant_city,omitempty"`
	TeacherContractions  string   `json:"teacher_contractions,omitempty"`
	TeacherTitle         string   `json:"teacher_title,omitempty"`
	Specialization       string   `json:"specialization,omitempty"`
}

// TenantSemester TODO: Add description
type TenantSemester struct {
	TenantID     int    `json:"tenant_id"`
	SemesterCode string `json:"semester_code"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	NPPCode      string `json:"npp_code"`
	// Additional fields
	NPPName      string `json:"npp_name,omitempty"`
	SemesterName string `json:"semester_name,omitempty"`
	FullName     string `json:"full_name,omitempty"`
}

// GetFullName TODO: Add description
func (t TenantSemester) GetFullName() string {
	return t.NPPName + " - " + t.SemesterName
}
