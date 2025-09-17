package wpmodels

type EnrollmentGrade struct {
	PupilID              int    `json:"pupil_id"`
	TenantID             int    `json:"tenant_id"`
	SubjectCode          string `json:"subject_code"`
	ClassCode            string `json:"class_code"`
	Grade                int    `json:"grade"`
	SchoolSpecialization string `json:"school_specialization"`
}
