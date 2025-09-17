package util

import (
	"database/sql"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"sort"
	"time"
)

// CreateLesson creates a new lesson record and associated pupil attendance records
// in a single database transaction. It first inserts the lesson data into the
// class_lesson table, retrieves the generated lesson ID, and then inserts all
// pupil attendance records for that lesson. If any operation fails, the entire
// transaction is rolled back to maintain data consistency.
//
// Parameters:
//   - requestData: Contains lesson details and attendance data for all pupils
//   - tenantDB: Database connection for the specific tenant
//
// Returns:
//   - error: nil on success, or error describing what went wrong
func CreateLesson(
	requestData tenantmodels.LessonData, tenantDB *sql.DB, signature string,
) (newLesson *tenantmodels.LessonData, err error) {
	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting tenantDB transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	lessonInsertQuery := `INSERT INTO class_lesson (description, date,
	period_number, section_id, subject_code, signature) VALUES (?, ?, ?, ?, ?, ?)`

	res, err := tx.Exec(
		lessonInsertQuery,
		requestData.LessonData.Description,
		requestData.LessonData.Date,
		requestData.LessonData.PeriodNumber,
		requestData.LessonData.SectionID,
		requestData.LessonData.SubjectCode,
		signature,
	)
	if err != nil {
		return nil, err
	}

	lessonID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	attendanceInsertQuery := `INSERT INTO pupil_attendance (pupil_id,
	lesson_id, status) VALUES (?, ?, ?)`

	for _, attendance := range requestData.PupilAttendanceData {
		_, err = tx.Exec(
			attendanceInsertQuery,
			attendance.PupilID,
			lessonID,
			attendance.Status,
		)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	newLesson, err = GetLessonByID(
		int(lessonID), tenantDB,
	)
	if err != nil {
		return nil, err
	}

	return newLesson, nil
}

// UpdateLesson updates an existing lesson record and its associated pupil attendance records
// in a single database transaction. It first updates the lesson data in the class_lesson table,
// then deletes all existing attendance records for that lesson and inserts the new ones.
// This approach ensures data consistency and handles cases where pupils may have been
// added or removed from the lesson. If any operation fails, the entire transaction is
// rolled back to maintain data integrity.
//
// Parameters:
//   - lessonID: The ID of the lesson to update
//   - requestData: Contains updated lesson details and attendance data for all pupils
//   - tenantDB: Database connection for the specific tenant
//
// Returns:
//   - error: nil on success, or error describing what went wrong
func UpdateLesson(
	lessonID int,
	requestData tenantmodels.LessonData,
	tenantDB *sql.DB,
	signature string,
) (updatedLesson *tenantmodels.LessonData, err error) {
	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting tenantDB transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Update lesson data
	lessonUpdateQuery := `UPDATE class_lesson SET description = ?, date = ?,
		period_number = ?, section_id = ?, subject_code = ?, signature = ? WHERE id = ?`

	_, err = tx.Exec(
		lessonUpdateQuery,
		requestData.LessonData.Description,
		requestData.LessonData.Date,
		requestData.LessonData.PeriodNumber,
		requestData.LessonData.SectionID,
		requestData.LessonData.SubjectCode,
		signature,
		lessonID,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating lesson: %v", err)
	}

	// Delete existing attendance records for this lesson
	deleteAttendanceQuery := `DELETE FROM pupil_attendance WHERE lesson_id = ?`
	_, err = tx.Exec(deleteAttendanceQuery, lessonID)
	if err != nil {
		return nil, fmt.Errorf("error deleting existing attendance records: %v", err)
	}

	// Insert new attendance records
	attendanceInsertQuery := `INSERT INTO pupil_attendance (pupil_id,
		lesson_id, status) VALUES (?, ?, ?)`

	for _, attendance := range requestData.PupilAttendanceData {
		_, err = tx.Exec(
			attendanceInsertQuery,
			attendance.PupilID,
			lessonID,
			attendance.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error inserting attendance record for pupil %d: %v",
				attendance.PupilID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedLesson, err = GetLessonByID(
		lessonID, tenantDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving updated lesson: %v", err)
	}

	return updatedLesson, nil
}

// DeleteLesson removes a lesson record from the database. The associated pupil
// attendance records are automatically deleted through database cascade constraints,
// ensuring data integrity and preventing orphaned attendance records.
//
// Parameters:
//   - lessonID: The ID of the lesson to delete
//   - tenantDB: Database connection for the specific tenant
//
// Returns:
//   - error: nil on success, or error describing what went wrong
func DeleteLesson(
	lessonID int,
	tenantDB *sql.DB,
) error {
	deleteLessonQuery := `DELETE FROM class_lesson WHERE id = ?`
	_, err := tenantDB.Exec(deleteLessonQuery, lessonID)
	if err != nil {
		return err
	}
	return nil
}

// GetLessonsForSection retrieves all lessons for a specific section along with
// their associated pupil attendance records. The function returns a slice of
// LessonData structs, each containing the lesson information and all attendance
// records for that lesson.
//
// Parameters:
//   - sectionID: The ID of the section to retrieve lessons for
//   - tenantDB: Database connection for the specific tenant
//
// Returns:
//   - []tenantmodels.LessonData: Slice of lessons with their attendance data
//   - error: nil on success, or error describing what went wrong
func GetLessonsForSection(
	sectionID int,
	tenantDB *sql.DB,
	claims *wpmodels.Claims,
) ([]tenantmodels.LessonData, error) {
	var rows *sql.Rows
	var err error
	var lessonQuery string
	// First, get all lessons for the section
	if claims.AccountType == "root" || claims.AccountType == "tenant_admin" {
		lessonQuery = `SELECT DISTINCT cl.id, cl.description, cl.date,
		cl.period_number, cl.section_id,
		cl.subject_code, s.subject_name, cl.signature
		FROM class_lesson cl
		JOIN ednevnik_workspace.subjects s
		ON s.subject_code = cl.subject_code
		WHERE cl.section_id = ? ORDER BY date DESC, s.subject_name ASC,
		period_number ASC`

		rows, err = tenantDB.Query(lessonQuery, sectionID)
	} else if claims.AccountType == "teacher" {
		lessonQuery = `SELECT DISTINCT cl.id, cl.description, cl.date,
		cl.period_number, cl.section_id,
		cl.subject_code, s.subject_name, cl.signature
		FROM class_lesson cl
		JOIN ednevnik_workspace.subjects s
		ON s.subject_code = cl.subject_code
		JOIN teachers_sections_subjects tss
		ON tss.subject_code = cl.subject_code
		WHERE cl.section_id = ? AND tss.teacher_id = ?
		ORDER BY date DESC, s.subject_name ASC, period_number ASC`

		rows, err = tenantDB.Query(lessonQuery, sectionID, claims.ID)
	} else {
		return nil, fmt.Errorf("unauthorized access for account type: %s", claims.AccountType)
	}
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error querying lessons: %v", err)
	}

	var lessons []tenantmodels.LessonData

	for rows.Next() {
		var lesson tenantmodels.ClassLesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Description,
			&lesson.Date,
			&lesson.PeriodNumber,
			&lesson.SectionID,
			&lesson.SubjectCode,
			&lesson.SubjectName,
			&lesson.Signature,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning lesson row: %v", err)
		}

		// Get attendance records for this lesson
		attendanceQuery := `SELECT pupil_id, lesson_id, status 
		FROM pupil_attendance WHERE lesson_id = ? ORDER BY pupil_id ASC`

		attendanceRows, err := tenantDB.Query(attendanceQuery, lesson.ID)
		if err != nil {
			return nil, fmt.Errorf("error querying attendance for lesson %d: %v", lesson.ID, err)
		}
		defer attendanceRows.Close()

		var attendanceRecords []tenantmodels.PupilAttendance
		for attendanceRows.Next() {
			var attendance tenantmodels.PupilAttendance
			err := attendanceRows.Scan(
				&attendance.PupilID,
				&attendance.LessonID,
				&attendance.Status,
			)
			if err != nil {
				return nil, fmt.Errorf("error scanning attendance row: %v", err)
			}
			attendanceRecords = append(attendanceRecords, attendance)
		}

		// Create LessonData struct with lesson and attendance data
		lessonData := tenantmodels.LessonData{
			LessonData:          lesson,
			PupilAttendanceData: attendanceRecords,
		}

		lessons = append(lessons, lessonData)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lesson rows: %v", err)
	}

	return lessons, nil
}

// GetLessonByID retrieves a specific lesson by its ID along with
// its associated pupil attendance records. The function returns a
// LessonData struct containing the lesson information and all attendance
// records for that lesson.
//
// Parameters:
//   - lessonID: The ID of the lesson to retrieve
//   - tenantDB: Database connection for the specific tenant
//
// Returns:
//   - *tenantmodels.LessonData: Pointer to lesson data with attendance records
//   - error: nil on success, or error describing what went wrong
func GetLessonByID(
	lessonID int,
	tenantDB *sql.DB,
) (*tenantmodels.LessonData, error) {
	// Get the lesson by ID
	lessonQuery := `SELECT cl.id, cl.description, cl.date, cl.period_number, cl.section_id,
	cl.subject_code, s.subject_name, cl.signature
	FROM class_lesson cl
	JOIN ednevnik_workspace.subjects s
	ON s.subject_code = cl.subject_code
	WHERE cl.id = ?`

	var lesson tenantmodels.ClassLesson
	err := tenantDB.QueryRow(lessonQuery, lessonID).Scan(
		&lesson.ID,
		&lesson.Description,
		&lesson.Date,
		&lesson.PeriodNumber,
		&lesson.SectionID,
		&lesson.SubjectCode,
		&lesson.SubjectName,
		&lesson.Signature,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lesson with ID %d not found", lessonID)
		}
		return nil, fmt.Errorf("error querying lesson: %v", err)
	}

	// Get attendance records for this lesson
	attendanceQuery := `SELECT pupil_id, lesson_id, status 
	FROM pupil_attendance WHERE lesson_id = ? ORDER BY pupil_id ASC`

	attendanceRows, err := tenantDB.Query(attendanceQuery, lessonID)
	if err != nil {
		return nil, fmt.Errorf("error querying attendance for lesson %d: %v", lessonID, err)
	}
	defer attendanceRows.Close()

	var attendanceRecords []tenantmodels.PupilAttendance
	for attendanceRows.Next() {
		var attendance tenantmodels.PupilAttendance
		err := attendanceRows.Scan(
			&attendance.PupilID,
			&attendance.LessonID,
			&attendance.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance row: %v", err)
		}
		attendanceRecords = append(attendanceRecords, attendance)
	}

	if err = attendanceRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance rows: %v", err)
	}

	// Create LessonData struct with lesson and attendance data
	lessonData := &tenantmodels.LessonData{
		LessonData:          lesson,
		PupilAttendanceData: attendanceRecords,
	}

	return lessonData, nil
}

// GetAbsentAttendancesForSectionHelper retrieves all absent, excused, and unexcused
// pupil attendance records for a specific section.
func GetAbsentAttendancesForSectionHelper(
	sectionID int,
	tenantDB *sql.DB,
) ([]tenantmodels.PupilAttendance, error) {
	query := `SELECT pa.pupil_id, pa.lesson_id, pa.status,
	p.name, p.last_name, cl.date, cl.period_number, s.subject_name
	FROM pupil_attendance pa
	JOIN class_lesson cl ON pa.lesson_id = cl.id
	JOIN pupils p ON pa.pupil_id = p.id
	JOIN ednevnik_workspace.subjects s ON cl.subject_code = s.subject_code
	WHERE cl.section_id = ?
	AND pa.status IN ('absent', 'excused', 'unexcused')`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying absent attendances: %v", err)
	}
	defer rows.Close()

	var attendances []tenantmodels.PupilAttendance
	for rows.Next() {
		var attendance tenantmodels.PupilAttendance
		err := rows.Scan(
			&attendance.PupilID,
			&attendance.LessonID,
			&attendance.Status,
			&attendance.Name,
			&attendance.LastName,
			&attendance.Date,
			&attendance.PeriodNumber,
			&attendance.SubjectName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance row: %v", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance rows: %v", err)
	}

	return attendances, nil
}

// GetAbsentAttendancesForPupilHelper retrieves all absent, excused, and unexcused
// pupil attendance records for a specific pupil within a section.
func GetAbsentAttendancesForPupilHelper(
	pupilID,
	sectionID int,
	tenantDB *sql.DB,
) ([]tenantmodels.PupilAttendance, error) {
	query := `SELECT pa.lesson_id, pa.status, cl.date, cl.period_number, s.subject_name
	FROM pupil_attendance pa
	JOIN class_lesson cl ON pa.lesson_id = cl.id
	JOIN ednevnik_workspace.subjects s ON cl.subject_code = s.subject_code
	WHERE pa.pupil_id = ? AND cl.section_id = ?
	AND pa.status IN ('absent', 'excused', 'unexcused')`

	rows, err := tenantDB.Query(query, pupilID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying absent attendances for pupil: %v", err)
	}
	defer rows.Close()

	var attendances []tenantmodels.PupilAttendance
	for rows.Next() {
		var attendance tenantmodels.PupilAttendance
		err := rows.Scan(
			&attendance.LessonID,
			&attendance.Status,
			&attendance.Date,
			&attendance.PeriodNumber,
			&attendance.SubjectName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attendance row: %v", err)
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attendance rows: %v", err)
	}

	return attendances, nil
}

// HandleAttendanceActionHelper processes attendance actions
func HandleAttendanceActionHelper(
	action tenantmodels.AttendanceAction, tenantDB *sql.DB,
) error {
	query := `UPDATE pupil_attendance SET status = ? WHERE pupil_id = ? AND lesson_id = ?`
	_, err := tenantDB.Exec(query, action.Type, action.PupilID, action.LessonID)
	return err
}

// WeekCountOfLessonsForSection returns the total number
// of lessons scheduled for a specific section.
func WeekCountOfLessonsForSection(
	sectionID int,
	tenantDB *sql.DB,
) (int, error) {
	query := `SELECT COUNT(*) FROM schedule WHERE section_id = ?`
	var count int
	err := tenantDB.QueryRow(query, sectionID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetLessonsByWeekForSection retrieves lessons grouped by week for a specific section.
func GetLessonsByWeekForSection(
	sectionID int,
	tenantDB *sql.DB,
) ([]tenantmodels.LessonWeekGroup, error) {
	query := `SELECT cl.id, cl.description, cl.date, cl.period_number, cl.section_id,
    cl.subject_code, s.subject_name, cl.signature
    FROM class_lesson cl
    JOIN ednevnik_workspace.subjects s ON s.subject_code = cl.subject_code
    WHERE cl.section_id = ?
    ORDER BY cl.date ASC, s.subject_name ASC, cl.period_number ASC`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying lessons: %v", err)
	}
	defer rows.Close()

	// Map to group lessons by week
	lessonsByWeek := make(map[string]map[string][]tenantmodels.ClassLesson)
	var weekOrder []string

	for rows.Next() {
		var lesson tenantmodels.ClassLesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.Description,
			&lesson.Date,
			&lesson.PeriodNumber,
			&lesson.SectionID,
			&lesson.SubjectCode,
			&lesson.SubjectName,
			&lesson.Signature,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning lesson row: %v", err)
		}

		// Parse the lesson date
		lessonTime, err := time.Parse("2006-01-02", lesson.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing lesson date %s: %v", lesson.Date, err)
		}

		// Calculate the Monday of the week for this lesson
		weekday := lessonTime.Weekday()
		daysFromMonday := int(weekday - time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7 // Handle Sunday (weekday 0)
		}
		mondayOfWeek := lessonTime.AddDate(0, 0, -daysFromMonday)

		// Calculate Sunday of the same week
		sundayOfWeek := mondayOfWeek.AddDate(0, 0, 6)

		// Format week string as "dd.mm.yyyy - dd.mm.yyyy"
		weekKey := fmt.Sprintf("%s - %s",
			mondayOfWeek.Format("02.01.2006"),
			sundayOfWeek.Format("02.01.2006"))

		// Initialize week map if it doesn't exist
		if _, exists := lessonsByWeek[weekKey]; !exists {
			lessonsByWeek[weekKey] = make(map[string][]tenantmodels.ClassLesson)
			weekOrder = append(weekOrder, weekKey)
		}

		// Group lessons by date within the week
		if _, exists := lessonsByWeek[weekKey][lesson.Date]; !exists {
			lessonsByWeek[weekKey][lesson.Date] = []tenantmodels.ClassLesson{}
		}
		lessonsByWeek[weekKey][lesson.Date] = append(lessonsByWeek[weekKey][lesson.Date], lesson)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lesson rows: %v", err)
	}

	// Get scheduled lessons count per week (constant for all weeks)
	scheduledPerWeek, err := WeekCountOfLessonsForSection(sectionID, tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error getting scheduled lessons count: %v", err)
	}

	// Build the result structure
	var result []tenantmodels.LessonWeekGroup
	for _, weekKey := range weekOrder {
		lessonDateMap := lessonsByWeek[weekKey]

		// Count held lessons (lessons that actually exist in class_lesson table)
		heldCount := 0
		for _, lessons := range lessonDateMap {
			heldCount += len(lessons)
		}

		// Calculate unheld lessons
		unheldCount := scheduledPerWeek - heldCount
		if unheldCount < 0 {
			unheldCount = 0 // In case there are more held lessons than scheduled
		}

		// Create ordered list of dates within the week
		var dateOrder []string
		for date := range lessonDateMap {
			dateOrder = append(dateOrder, date)
		}

		// Sort dates within the week
		sort.Strings(dateOrder)

		// Build LessonDateGroup slice for this week
		var lessonsByDate []tenantmodels.LessonDateGroup
		for _, date := range dateOrder {
			lessonsByDate = append(lessonsByDate, tenantmodels.LessonDateGroup{
				Date:    date,
				Lessons: lessonDateMap[date],
			})
		}

		result = append(result, tenantmodels.LessonWeekGroup{
			Week:                weekKey,
			LessonsByDate:       lessonsByDate,
			TotalLessonsInWeek:  scheduledPerWeek,
			HeldLessonsInWeek:   heldCount,
			UnheldLessonsInWeek: unheldCount,
		})
	}

	return result, nil
}

func GetAbsencesByWeekForSection(
	sectionID int,
	tenantDB *sql.DB,
) ([]tenantmodels.WeekAbsenceGroup, error) {
	query := `SELECT pa.pupil_id, p.name, p.last_name, pa.lesson_id, pa.status,
        cl.date, cl.period_number, cl.subject_code, s.subject_name
        FROM pupil_attendance pa
        JOIN class_lesson cl ON pa.lesson_id = cl.id
        JOIN pupils p ON pa.pupil_id = p.id
        JOIN ednevnik_workspace.subjects s ON cl.subject_code = s.subject_code
        WHERE cl.section_id = ? AND pa.status IN ('absent', 'excused', 'unexcused')
        ORDER BY cl.date ASC, s.subject_name ASC, p.last_name ASC`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying absences: %v", err)
	}
	defer rows.Close()

	// structure: week -> date -> subject_code -> []PupilAttendance
	absencesByWeek := make(map[string]map[string]map[string][]tenantmodels.PupilAttendance)
	var weekOrder []string

	for rows.Next() {
		var rec tenantmodels.PupilAttendance
		if err := rows.Scan(
			&rec.PupilID,
			&rec.Name,
			&rec.LastName,
			&rec.LessonID,
			&rec.Status,
			&rec.Date,
			&rec.PeriodNumber,
			&rec.SubjectCode,
			&rec.SubjectName,
		); err != nil {
			return nil, fmt.Errorf("error scanning absence row: %v", err)
		}

		lessonTime, err := time.Parse("2006-01-02", rec.Date)
		if err != nil {
			return nil, fmt.Errorf("error parsing lesson date %s: %v", rec.Date, err)
		}

		weekday := lessonTime.Weekday()
		daysFromMonday := int(weekday - time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7
		}
		mondayOfWeek := lessonTime.AddDate(0, 0, -daysFromMonday)
		sundayOfWeek := mondayOfWeek.AddDate(0, 0, 6)
		weekKey := fmt.Sprintf("%s - %s",
			mondayOfWeek.Format("02.01.2006"),
			sundayOfWeek.Format("02.01.2006"),
		)

		if _, ok := absencesByWeek[weekKey]; !ok {
			absencesByWeek[weekKey] = make(map[string]map[string][]tenantmodels.PupilAttendance)
			weekOrder = append(weekOrder, weekKey)
		}
		if _, ok := absencesByWeek[weekKey][rec.Date]; !ok {
			absencesByWeek[weekKey][rec.Date] = make(map[string][]tenantmodels.PupilAttendance)
		}
		absencesByWeek[weekKey][rec.Date][rec.SubjectCode] = append(absencesByWeek[weekKey][rec.Date][rec.SubjectCode], rec)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating absence rows: %v", err)
	}

	// Build result
	var result []tenantmodels.WeekAbsenceGroup
	for _, wk := range weekOrder {
		dateMap := absencesByWeek[wk]

		// Collect and sort dates
		var dateOrder []string
		for d := range dateMap {
			dateOrder = append(dateOrder, d)
		}
		sort.Strings(dateOrder)

		var dayGroups []tenantmodels.DayAbsenceGroup
		weekTotal := 0
		weekExc := 0
		weekUnexc := 0
		weekPend := 0

		for _, d := range dateOrder {
			subjMap := dateMap[d]

			// sort subject codes
			var subjOrder []string
			for sc := range subjMap {
				subjOrder = append(subjOrder, sc)
			}
			sort.Strings(subjOrder)

			var daySubjects []tenantmodels.DaySubjectAbsenceGroup
			dayTotal := 0
			dayExc := 0
			dayUnexc := 0
			dayPend := 0

			for _, sc := range subjOrder {
				list := subjMap[sc]
				subjTotal := len(list)
				subjExc := 0
				subjUnexc := 0
				subjPend := 0
				for _, a := range list {
					switch a.Status {
					case "excused":
						subjExc++
					case "unexcused":
						subjUnexc++
					default:
						subjPend++
					}
				}
				daySubjects = append(daySubjects, tenantmodels.DaySubjectAbsenceGroup{
					SubjectCode: sc,
					SubjectName: "", // popuniti iz list[0] ako treba
					Absences:    list,
					Total:       subjTotal,
					Excused:     subjExc,
					Unexcused:   subjUnexc,
					Pending:     subjPend,
				})
				if len(list) > 0 && daySubjects[len(daySubjects)-1].SubjectName == "" {
					daySubjects[len(daySubjects)-1].SubjectName = list[0].SubjectName
				}
				dayTotal += subjTotal
				dayExc += subjExc
				dayUnexc += subjUnexc
				dayPend += subjPend
			}

			dayGroups = append(dayGroups, tenantmodels.DayAbsenceGroup{
				Date:      d,
				Subjects:  daySubjects,
				Total:     dayTotal,
				Excused:   dayExc,
				Unexcused: dayUnexc,
				Pending:   dayPend,
			})

			weekTotal += dayTotal
			weekExc += dayExc
			weekUnexc += dayUnexc
			weekPend += dayPend
		}

		result = append(result, tenantmodels.WeekAbsenceGroup{
			Week:      wk,
			Days:      dayGroups,
			Total:     weekTotal,
			Excused:   weekExc,
			Unexcused: weekUnexc,
			Pending:   weekPend,
		})
	}

	return result, nil
}
