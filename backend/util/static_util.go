package util

import (
	"database/sql"
	"ednevnik-backend/models/interfaces"
	wpmodels "ednevnik-backend/models/workspace"
)

// InsertCanton TODO: Add description
func InsertCanton(tx *sql.Tx, canton wpmodels.Canton) error {
	query := "INSERT IGNORE INTO cantons (canton_code, canton_name, country) VALUES (?, ?, ?)"
	_, err := tx.Exec(query, canton.CantonCode, canton.CantonName, canton.Country)
	return err
}

// InsertClass TODO: Add description
func InsertClass(tx *sql.Tx, class wpmodels.Class) error {
	query := "INSERT IGNORE INTO classes (class_code) VALUES (?)"
	_, err := tx.Exec(query, class.ClassCode)
	return err
}

// InsertSubject TODO: Add description
func InsertSubject(tx *sql.Tx, subject wpmodels.Subject) error {
	query := "INSERT IGNORE INTO subjects (subject_name, subject_code) VALUES (?, ?)"
	_, err := tx.Exec(query, subject.SubjectName, subject.SubjectCode)
	return err
}

// InsertNPP TODO: Add description
func InsertNPP(tx *sql.Tx, npp wpmodels.NPP) error {
	query := "INSERT IGNORE INTO npp (npp_name, npp_code) VALUES (?, ?)"
	_, err := tx.Exec(query, npp.NPPName, npp.NPPCode)
	return err
}

// InsertCourse TODO: Add description
func InsertCourse(tx *sql.Tx, course wpmodels.Course) error {
	query := "INSERT IGNORE INTO courses_secondary (course_code, course_name, course_duration) VALUES (?, ?, ?)"
	_, err := tx.Exec(query, course.CourseCode, course.CourseName, course.CourseDuration)
	return err
}

// InsertCurriculum TODO: Add description
func InsertCurriculum(tx *sql.Tx, curriculum wpmodels.CurriculumCreate) error {
	query := `INSERT IGNORE INTO curriculum (curriculum_code, curriculum_name,
	class_code, npp_code, course_code, canton_code, tenant_type, final_curriculum)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		curriculum.CurriculumCode,
		curriculum.CurriculumName,
		curriculum.ClassCode,
		curriculum.NPPCode,
		curriculum.CourseCode,
		curriculum.CantonCode,
		curriculum.TenantType,
		curriculum.FinalCurriculum,
	)
	return err
}

// GetCurriculumByCode TODO: Add description
func GetCurriculumByCode(db interfaces.DatabaseQuerier, curriculumCode string) (wpmodels.Curriculum, error) {
	var curriculum wpmodels.Curriculum
	query := `SELECT curriculum_code, curriculum_name, class_code, npp_code,
	course_code, canton_code, tenant_type FROM curriculum WHERE curriculum_code = ?`
	err := db.QueryRow(query, curriculumCode).Scan(
		&curriculum.CurriculumCode,
		&curriculum.CurriculumName,
		&curriculum.ClassCode,
		&curriculum.NPPCode,
		&curriculum.CourseCode,
		&curriculum.CantonCode,
		&curriculum.TenantType,
	)
	if err != nil {
		return curriculum, err
	}
	return curriculum, nil
}

// InsertCurriculumSubject TODO: Add description
func InsertCurriculumSubject(tx *sql.Tx, curriculumSubject wpmodels.CurriculumSubject) error {
	query := "INSERT IGNORE INTO curriculum_subjects (curriculum_code, subject_code) VALUES (?, ?)"
	_, err := tx.Exec(query, curriculumSubject.CurriculumCode, curriculumSubject.SubjectCode)
	return err
}

// InsertSemester TODO: Add description
func InsertSemester(tx *sql.Tx, semester wpmodels.Semester) error {
	query := "INSERT IGNORE INTO semester (semester_code, semester_name, progress_level) VALUES (?, ?, ?)"
	_, err := tx.Exec(query, semester.SemesterCode, semester.SemesterName, semester.ProgressLevel)
	return err
}

// GetSemesterByCode TODO: Add description
func GetSemesterByCode(db *sql.DB, semesterCode string) (wpmodels.Semester, error) {
	var semester wpmodels.Semester
	query := `SELECT semester_code, semester_name, progress_level FROM
	semester WHERE semester_code = ?`
	err := db.QueryRow(query, semesterCode).Scan(
		&semester.SemesterCode,
		&semester.SemesterName,
		&semester.ProgressLevel,
	)
	if err != nil {
		return semester, err
	}
	return semester, nil
}

// InsertNPPSemester TODO: Add description
func InsertNPPSemester(tx *sql.Tx, nppSemester wpmodels.NPPSemester) error {
	query := `INSERT IGNORE INTO npp_semester (npp_code, semester_code,
	start_date, end_date) VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(query, nppSemester.NPPCode, nppSemester.SemesterCode,
		nppSemester.StartDate, nppSemester.EndDate)
	return err
}

// GetSemestersByNPPCode TODO: Add description
func GetSemestersByNPPCode(db *sql.DB, nppCode string) ([]wpmodels.NPPSemester, error) {
	query := `SELECT npp_code, semester_code, start_date, end_date
	FROM npp_semester WHERE npp_code = ?`
	rows, err := db.Query(query, nppCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var NPPSSemesters []wpmodels.NPPSemester
	for rows.Next() {
		var nppSemester wpmodels.NPPSemester
		if err := rows.Scan(&nppSemester.NPPCode, &nppSemester.SemesterCode,
			&nppSemester.StartDate, &nppSemester.EndDate); err != nil {
			return nil, err
		}
		NPPSSemesters = append(NPPSSemesters, nppSemester)
	}
	return NPPSSemesters, nil
}

// GetAllSubjectsForCurriculumCode TODO: Add description
func GetAllSubjectsForCurriculumCode(
	curriculumCode string, workspaceDB interfaces.DatabaseQuerier,
) ([]wpmodels.Subject, error) {
	query := `SELECT s.subject_name, s.subject_code FROM subjects s
	JOIN curriculum_subjects cs ON s.subject_code = cs.subject_code
	WHERE cs.curriculum_code = ? ORDER BY s.subject_name`
	rows, err := workspaceDB.Query(query, curriculumCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []wpmodels.Subject
	for rows.Next() {
		var subject wpmodels.Subject
		if err := rows.Scan(&subject.SubjectName, &subject.SubjectCode); err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}
