package wpmodels

// Canton TODO: Add description
type Canton struct {
	CantonCode string `json:"canton_code"`
	CantonName string `json:"canton_name"`
	Country    string `json:"country"`
}

// Class TODO: Add description
type Class struct {
	ClassCode string `json:"class_code"`
}

// Subject TODO: Add description
type Subject struct {
	SubjectCode string `json:"subject_code"`
	SubjectName string `json:"subject_name"`
	// Optional field for handling teacher subject invites
	Checked       bool `json:"checked,omitempty"`
	InviteIndexID *int `json:"invite_index_id,omitempty"`
}

// NPP TODO: Add description
type NPP struct {
	NPPCode string `json:"npp_code"`
	NPPName string `json:"npp_name"`
}

// Course TODO: Add description
type Course struct {
	CourseCode     string `json:"course_code"`
	CourseName     string `json:"course_name"`
	CourseDuration string `json:"course_duration"`
}

// Curriculum TODO: Add description
type Curriculum struct {
	CurriculumCode string  `json:"curriculum_code"`
	CurriculumName string  `json:"curriculum_name"`
	ClassCode      string  `json:"class_code"`
	NPPName        string  `json:"npp_name,omitempty"`
	NPPCode        string  `json:"npp_code"`
	CourseCode     *string `json:"course_code,omitempty"`
	CantonCode     string  `json:"canton_code"`
	TenantType     string  `json:"tenant_type"`
	CourseName     string  `json:"course_name,omitempty"`
}

type CurriculumCreate struct {
	CurriculumCode  string  `json:"curriculum_code"`
	CurriculumName  string  `json:"curriculum_name"`
	ClassCode       string  `json:"class_code"`
	NPPName         string  `json:"npp_name,omitempty"`
	NPPCode         string  `json:"npp_code"`
	CourseCode      *string `json:"course_code,omitempty"`
	CantonCode      string  `json:"canton_code"`
	TenantType      string  `json:"tenant_type"`
	FinalCurriculum bool    `json:"final_curriculum"`
}

// CurriculumGet TODO: Add description
type CurriculumGet struct {
	CurriculumCode string  `json:"curriculum_code"`
	CurriculumName string  `json:"curriculum_name"`
	ClassCode      string  `json:"class_code"`
	NPPName        string  `json:"npp_name"`
	CourseName     *string `json:"course_name,omitempty"`
	CantonCode     string  `json:"canton_code"`
	TenantType     string  `json:"tenant_type"`
}

// CurriculumSubject TODO: Add description
type CurriculumSubject struct {
	CurriculumCode string `json:"curriculum_code"`
	SubjectCode    string `json:"subject_code"`
}

// Semester TODO: Add description
type Semester struct {
	SemesterCode  string `json:"semester_code"`
	SemesterName  string `json:"semester_name"`
	ProgressLevel int    `json:"progress_level"`
}

// NPPSemester TODO: Add description
type NPPSemester struct {
	NPPCode      string `json:"npp_code"`
	SemesterCode string `json:"semester_code"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	// Additional fields
	NPPName      string `json:"npp_name,omitempty"`
	SemesterName string `json:"semester_name,omitempty"`
	FullName     string `json:"full_name,omitempty"`
}

// GetFullName TODO: Add description
func (s NPPSemester) GetFullName() string {
	return s.NPPName + " - " + s.SemesterName
}
