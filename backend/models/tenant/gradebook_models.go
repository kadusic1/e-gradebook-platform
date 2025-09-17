package tenantmodels

import wpmodels "ednevnik-backend/models/workspace"

// Grade represents a grade for a pupil
type Grade struct {
	ID          int    `json:"id"`
	PupilID     int    `json:"pupil_id"`
	SectionID   int    `json:"section_id"`
	SubjectCode string `json:"subject_code"`
	Grade       int    `json:"grade"`
	GradeDate   string `json:"grade_date"`
	TeacherID   int    `json:"teacher_id,omitempty"`
	Type        string `json:"type"`
	Signature   string `json:"signature"`
	IsEdited    bool   `json:"is_edited,omitempty"`
	IsDeleted   bool   `json:"is_deleted,omitempty"`
	ValidUntil  string `json:"valid_until,omitempty"`
	// Optional fields
	TeacherName     string `json:"teacher_name,omitempty"`
	TeacherLastName string `json:"teacher_last_name,omitempty"`
	SubjectName     string `json:"subject_name,omitempty"`
	SemesterCode    string `json:"semester_code,omitempty"`
}

// GradePupilGroup represents a group of grades for a pupil
type GradePupilGroup struct {
	Pupil        Pupil   `json:"pupil"`
	AverageGrade float64 `json:"average_grade,omitempty"`
	Grades       []Grade `json:"grades"`
}

// BehaviourGrade representsa a pupil behaviour grade
type BehaviourGrade struct {
	ID           int    `json:"id"`
	PupilID      int    `json:"pupil_id"`
	SectionID    int    `json:"section_id"`
	Behaviour    string `json:"behaviour"`
	SemesterCode string `json:"semester_code"`
	Signature    string `json:"behaviour_determined_by_teacher"`
	ValidUntil   string `json:"valid_until,omitempty"`
	// Pupil name and last name
	Name         string `json:"name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	SemesterName string `json:"semester_name,omitempty"`
	Date         string `json:"date,omitempty"`
}

// CompleteGradebook contains all gradebook data including pupils, schedules, grades,
// behaviour grades and lessons
type CompleteGradebook struct {
	Pupils          []Pupil                   `json:"pupils"`
	ScheduleHistory []ScheduleGroupCollection `json:"schedule_history"`
	GradeData       []CompleteGradebookData   `json:"grade_data"`
	Subjects        []wpmodels.Subject        `json:"subjects"`
	Lessons         []LessonWeekGroup         `json:"lessons"`
	Absences        []WeekAbsenceGroup        `json:"absences"`
}

// SubjectGradeGroup represents a group of grades for a specific subject
type SubjectGradeGroup struct {
	SubjectName string  `json:"subject_name"`
	SubjectCode string  `json:"subject_code"`
	Grades      []Grade `json:"grades"`
}

// SemesterGradeGroup represents a group of grades for a specific semester
type SemesterGradeGroup struct {
	SemesterName    string              `json:"semester_name"`
	SubjectGrades   []SubjectGradeGroup `json:"subject_grades"`
	BehaviourGrades []BehaviourGrade    `json:"behaviour_grades"`
}

// CompleteGradebookData represents the complete gradebook data for a pupil
type CompleteGradebookData struct {
	PupilName         string               `json:"pupil_name"`
	PupilUnenrolled   bool                 `json:"pupil_unenrolled"`
	GradesForSemester []SemesterGradeGroup `json:"grades_for_semester"`
}
