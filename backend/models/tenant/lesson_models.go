package tenantmodels

// ClassLesson is a struct containing fields for lesson info (školski čas)
type ClassLesson struct {
	ID           int    `json:"id,omitempty"`
	Description  string `json:"description"`
	Date         string `json:"date,omitempty"`
	PeriodNumber int    `json:"period_number"`
	SectionID    int    `json:"section_id"`
	SubjectCode  string `json:"subject_code"`
	SubjectName  string `json:"subject_name,omitempty"`
	Signature    string `json:"lesson_posted_by_teacher"`
}

// PupilAttendance is a struct containing fields regarding attendance.
// Attendance status can be: present, absent, unexcused, excused.
type PupilAttendance struct {
	PupilID      int    `json:"pupil_id"`
	LessonID     int    `json:"lesson_id,omitempty"`
	Status       string `json:"status"`
	Name         string `json:"name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Date         string `json:"date,omitempty"`
	PeriodNumber int    `json:"period_number,omitempty"`
	SubjectName  string `json:"subject_name,omitempty"`
	SubjectCode  string `json:"subject_code,omitempty"`
}

// LessonData is a struct used for creating or updating a lesson along with
// the attendance data for all pupils in that lesson. It combines lesson details
// with the corresponding attendance records in a single request payload.
type LessonData struct {
	LessonData          ClassLesson       `json:"lesson_data"`
	PupilAttendanceData []PupilAttendance `json:"pupil_attendance_data"`
}

// AttendanceAction is a struct used to represent an action on a pupil's attendance.
type AttendanceAction struct {
	Type     string `json:"type"`
	PupilID  int    `json:"pupil_id"`
	LessonID int    `json:"lesson_id"`
}

// LessonDateGroup is a struct used to group lessons by date.
type LessonDateGroup struct {
	Date    string        `json:"date"`
	Lessons []ClassLesson `json:"lessons_for_date"`
}

// LessonWeekGroup is a struct used to group lessons by week.
type LessonWeekGroup struct {
	Week                string            `json:"week"`
	LessonsByDate       []LessonDateGroup `json:"lesson_date_group"`
	TotalLessonsInWeek  int               `json:"total_lessons_in_week"`
	HeldLessonsInWeek   int               `json:"held_lessons_in_week"`
	UnheldLessonsInWeek int               `json:"unheld_lessons_in_week"`
}

// DaySubjectAbsenceGroup groups absences for a single subject on a single date.
type DaySubjectAbsenceGroup struct {
	SubjectCode string            `json:"subject_code"`
	SubjectName string            `json:"subject_name"`
	Absences    []PupilAttendance `json:"absences"`
	Total       int               `json:"total"`
	Excused     int               `json:"excused_count"`
	Unexcused   int               `json:"unexcused_count"`
	Pending     int               `json:"pending_count"` // status 'absent' / not yet marked
}

// DayAbsenceGroup groups absences for a single date, organized by subject.
type DayAbsenceGroup struct {
	Date      string                   `json:"date"`
	Subjects  []DaySubjectAbsenceGroup `json:"subjects"`
	Total     int                      `json:"total"`
	Excused   int                      `json:"excused_count"`
	Unexcused int                      `json:"unexcused_count"`
	Pending   int                      `json:"pending_count"`
}

// WeekAbsenceGroup predstavlja sve izostanke za jednu sedmicu, najprije razvrstane po danima, pa po predmetima.
type WeekAbsenceGroup struct {
	Week      string            `json:"week"`
	Days      []DayAbsenceGroup `json:"days"`
	Total     int               `json:"total"`
	Excused   int               `json:"excused_count"`
	Unexcused int               `json:"unexcused_count"`
	Pending   int               `json:"pending_count"`
}
