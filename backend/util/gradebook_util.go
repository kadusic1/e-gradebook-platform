package util

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	interfaces "ednevnik-backend/models/interfaces"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"math"
)

// CalculateAverageGrade calculates the average grade from a slice of grades
// Returns the result rounded to exactly 2 decimal places
func CalculateAverageGrade(grades []tenantmodels.Grade) float64 {
	if len(grades) == 0 {
		return 0.0
	}

	var total float64
	var count int

	for _, grade := range grades {
		if grade.Type == "final" || grade.IsDeleted {
			continue
		}
		total += float64(grade.Grade)
		count++
	}

	// If no non-final grades found, return 0
	if count == 0 {
		return 0.0
	}

	average := total / float64(count)

	// Round to 2 decimal places
	return math.Round(average*100) / 100
}

// CalculateAverageFinalGrade calculates the average grade from a slice of
// final grades. Returns the result rounded to exactly 2 decimal places
func CalculateAverageFinalGrade(grades []tenantmodels.Grade) float64 {
	if len(grades) == 0 {
		return 0.0
	}

	var total float64
	var count int

	for _, grade := range grades {
		if grade.Type != "final" {
			continue
		}
		total += float64(grade.Grade)
		count++
	}

	// If no non-final grades found, return 0
	if count == 0 {
		return 0.0
	}

	average := total / float64(count)

	// Round to 2 decimal places
	return math.Round(average*100) / 100
}

// GetGradesForSectionSubject retrieves the grades in a specific section
// and for a specific subject, including deleted grades.
func GetGradesForSectionSubject(
	sectionID int,
	semesterCode,
	subjectCode string,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Grade, error) {
	gradesQuery := `
	WITH all_grades AS (
		SELECT sg.id, sg.pupil_id, sg.section_id, sg.subject_code, sg.grade,
		sg.grade_date, sg.type, sg.teacher_id, sg.semester_code, sg.signature,
		sg.ROW_START, sg.ROW_END
		FROM student_grades FOR SYSTEM_TIME ALL sg
		WHERE sg.section_id = ? AND sg.subject_code = ?
		AND sg.semester_code = ?
	),
	current_grade_ids AS (
		SELECT id FROM student_grades
		WHERE section_id = ? AND subject_code = ? AND semester_code = ?
	),
	latest_versions AS (
		SELECT id, MAX(ROW_START) as latest_row_start
		FROM all_grades GROUP BY id
	)
	SELECT ag.id, ag.pupil_id, ag.section_id, ag.subject_code, ag.grade,
	ag.grade_date, ag.type, COALESCE(t.name, '') as teacher_name,
	COALESCE(t.last_name, '') as teacher_last_name, s.subject_name,
	ag.semester_code, ag.signature,
	CASE
		WHEN ag.id IN (SELECT id FROM current_grade_ids)
			AND EXISTS (
				SELECT 1 FROM all_grades hist
				WHERE hist.id = ag.id
				AND hist.ROW_START < ag.ROW_START
			)
		THEN 1
		ELSE 0
	END AS is_edited,
	CASE
		WHEN ag.id NOT IN (SELECT id FROM current_grade_ids) THEN 1
		ELSE 0
	END AS is_deleted
	FROM all_grades ag
	JOIN latest_versions lv ON ag.id = lv.id AND ag.ROW_START = lv.latest_row_start
	JOIN ednevnik_workspace.subjects s ON ag.subject_code = s.subject_code
	LEFT JOIN ednevnik_workspace.teachers t ON ag.teacher_id = t.id
	ORDER BY is_deleted ASC, ag.grade_date ASC`

	rows, err := tenantDB.Query(gradesQuery,
		sectionID, subjectCode, semesterCode, // for all_grades CTE
		sectionID, subjectCode, semesterCode, // for current_grade_ids CTE
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []tenantmodels.Grade
	for rows.Next() {
		var grade tenantmodels.Grade
		var isEdited, isDeleted int

		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.Type,
			&grade.TeacherName, &grade.TeacherLastName, &grade.SubjectName,
			&grade.SemesterCode, &grade.Signature, &isEdited, &isDeleted,
		); err != nil {
			return nil, err
		}

		// Convert int to bool
		grade.IsEdited = isEdited == 1
		grade.IsDeleted = isDeleted == 1

		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

// GetGradesForSectionPupil retrieves the grades in a specific section
// and for a specific pupil, including deleted grades.
func GetGradesForSectionPupil(
	sectionID,
	pupilID int,
	semesterCode string,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Grade, error) {
	gradesQuery := `
	WITH all_grades AS (
		SELECT sg.id, sg.pupil_id, sg.section_id, sg.subject_code, sg.grade,
		sg.grade_date, sg.type, sg.teacher_id, sg.semester_code, sg.signature,
		sg.ROW_START, sg.ROW_END
		FROM student_grades FOR SYSTEM_TIME ALL sg
		WHERE sg.pupil_id = ? AND sg.section_id = ? AND sg.semester_code = ?
	),
	current_grade_ids AS (
		SELECT id FROM student_grades
		WHERE pupil_id = ? AND section_id = ? AND semester_code = ?
	),
	latest_versions AS (
		SELECT id, MAX(ROW_START) as latest_row_start
		FROM all_grades GROUP BY id
	)
	SELECT ag.id, ag.pupil_id, ag.section_id, ag.subject_code, ag.grade,
	ag.grade_date, ag.type, COALESCE(t.name, '') as teacher_name,
	COALESCE(t.last_name, '') as teacher_last_name, s.subject_name,
	ag.semester_code, ag.signature,
	CASE
		WHEN ag.id IN (SELECT id FROM current_grade_ids)
			AND EXISTS (
				SELECT 1 FROM all_grades hist
				WHERE hist.id = ag.id
				AND hist.ROW_START < ag.ROW_START
			)
		THEN 1
		ELSE 0
	END AS is_edited,
	CASE
		WHEN ag.id NOT IN (SELECT id FROM current_grade_ids) THEN 1
		ELSE 0
	END AS is_deleted
	FROM all_grades ag
	JOIN latest_versions lv ON ag.id = lv.id AND ag.ROW_START = lv.latest_row_start
	JOIN ednevnik_workspace.subjects s ON ag.subject_code = s.subject_code
	LEFT JOIN ednevnik_workspace.teachers t ON ag.teacher_id = t.id
	ORDER BY is_deleted ASC, ag.grade_date ASC`

	rows, err := tenantDB.Query(gradesQuery,
		pupilID, sectionID, semesterCode, // for all_grades CTE
		pupilID, sectionID, semesterCode, // for current_grade_ids CTE
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []tenantmodels.Grade
	for rows.Next() {
		var grade tenantmodels.Grade
		var isEdited, isDeleted int

		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.Type,
			&grade.TeacherName, &grade.TeacherLastName, &grade.SubjectName,
			&grade.SemesterCode, &grade.Signature, &isEdited, &isDeleted,
		); err != nil {
			return nil, err
		}

		// Convert int to bool
		grade.IsEdited = isEdited == 1
		grade.IsDeleted = isDeleted == 1

		grades = append(grades, grade)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

// GetPupilGradesForSectionPupilHelper retrieves all grades for a pupil in a section
func GetPupilGradesForSectionPupilHelper(
	sectionID,
	pupilID int,
	semesterCode string,
	tenantDB interfaces.DatabaseQuerier,
	workspaceDB interfaces.DatabaseQuerier,
) ([]commonmodels.GradeSubjectGroup, error) {

	section, err := GetSectionByID(
		int64(sectionID),
		tenantDB,
	)
	if err != nil {
		return nil, err
	}

	subjects, err := GetAllSubjectsForCurriculumCode(
		section.CurriculumCode, workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	grades, err := GetGradesForSectionPupil(
		sectionID,
		pupilID,
		semesterCode,
		tenantDB,
	)
	if err != nil {
		return nil, err
	}

	var gradeSubjectGroups []commonmodels.GradeSubjectGroup
	for _, subject := range subjects {
		subjectGrades := findGradesForSubject(subject.SubjectCode, grades)
		gradeSubjectGroups = append(gradeSubjectGroups, commonmodels.GradeSubjectGroup{
			Subject:      subject,
			Grades:       subjectGrades,
			AverageGrade: CalculateAverageGrade(subjectGrades),
		})
	}

	return gradeSubjectGroups, nil
}

func findGradesForPupil(
	pupilID int,
	grades []tenantmodels.Grade,
) []tenantmodels.Grade {
	var pupilGrades []tenantmodels.Grade
	for _, grade := range grades {
		if grade.PupilID == pupilID {
			pupilGrades = append(pupilGrades, grade)
		}
	}
	return pupilGrades
}

func findGradesForSubject(
	subjectCode string,
	grades []tenantmodels.Grade,
) []tenantmodels.Grade {
	var subjectGrades []tenantmodels.Grade
	for _, grade := range grades {
		if grade.SubjectCode == subjectCode {
			subjectGrades = append(subjectGrades, grade)
		}
	}
	return subjectGrades
}

// GetPupilGradesForSectionSubject retrieves the grades for all pupils in a section
func GetPupilGradesForSectionSubject(
	sectionID int,
	semesterCode,
	subjectCode string,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.GradePupilGroup, error) {
	pupils, err := GetPupilsForSection(
		fmt.Sprintf("%d", sectionID),
		true,
		tenantDB,
	)
	if err != nil {
		return nil, err
	}

	grades, err := GetGradesForSectionSubject(
		sectionID,
		semesterCode,
		subjectCode,
		tenantDB,
	)
	if err != nil {
		return nil, err
	}

	var gradePupilsGroup []tenantmodels.GradePupilGroup
	for _, pupil := range pupils {
		pupilGrades := findGradesForPupil(pupil.ID, grades)
		gradePupilsGroup = append(gradePupilsGroup, tenantmodels.GradePupilGroup{
			Pupil:        pupil,
			Grades:       pupilGrades,
			AverageGrade: CalculateAverageGrade(pupilGrades),
		})
	}

	return gradePupilsGroup, nil
}

// GetGradesForSectionSubjectPupil retrieves the grades in a specific section,
// specific subject and a specific pupil, including deleted grades.
func GetGradesForSectionSubjectPupil(
	sectionID int,
	subjectCode string,
	pupilID int,
	semesterCode string,
	tenantDB interfaces.DatabaseQuerier,
	workspaceDB interfaces.DatabaseQuerier,
) (*tenantmodels.GradePupilGroup, error) {
	pupil, err := GetGlobalPupilByID(
		fmt.Sprintf("%d", pupilID),
		workspaceDB,
	)
	pupil.Password = "" // Clear password for security
	if err != nil {
		return nil, err
	}

	gradesQuery := `
	WITH all_grades AS (
		SELECT sg.id, sg.pupil_id, sg.section_id, sg.subject_code, sg.grade,
		sg.grade_date, sg.type, sg.teacher_id, sg.semester_code, sg.signature,
		sg.ROW_START, sg.ROW_END
		FROM student_grades FOR SYSTEM_TIME ALL sg
		WHERE sg.pupil_id = ? AND sg.section_id = ? AND sg.subject_code = ?
		AND sg.semester_code = ?
	),
	current_grade_ids AS (
		SELECT id FROM student_grades
		WHERE pupil_id = ? AND section_id = ? AND subject_code = ?
		AND semester_code = ?
	),
	latest_versions AS (
		SELECT id, MAX(ROW_START) as latest_row_start
		FROM all_grades GROUP BY id
	)
	SELECT ag.id, ag.pupil_id, ag.section_id, ag.subject_code, ag.grade,
	ag.grade_date, ag.type, COALESCE(t.name, '') as teacher_name,
	COALESCE(t.last_name, '') as teacher_last_name, s.subject_name,
	ag.semester_code, ag.signature,
	CASE
		WHEN ag.id IN (SELECT id FROM current_grade_ids)
			AND EXISTS (
				SELECT 1 FROM all_grades hist
				WHERE hist.id = ag.id
				AND hist.ROW_START < ag.ROW_START
			)
		THEN 1
		ELSE 0
	END AS is_edited,
	CASE
		WHEN ag.id NOT IN (SELECT id FROM current_grade_ids) THEN 1
		ELSE 0
	END AS is_deleted
	FROM all_grades ag
	JOIN latest_versions lv ON ag.id = lv.id AND ag.ROW_START = lv.latest_row_start
	JOIN pupils p ON ag.pupil_id = p.id
	JOIN ednevnik_workspace.subjects s ON ag.subject_code = s.subject_code
	LEFT JOIN ednevnik_workspace.teachers t ON ag.teacher_id = t.id
	ORDER BY is_deleted ASC, ag.grade_date ASC`

	rows, err := tenantDB.Query(
		gradesQuery,
		pupilID, sectionID, subjectCode, semesterCode, // for all_grades CTE
		pupilID, sectionID, subjectCode, semesterCode, // for current_grade_ids CTE
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []tenantmodels.Grade
	for rows.Next() {
		var grade tenantmodels.Grade
		var isEdited, isDeleted int

		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.Type,
			&grade.TeacherName, &grade.TeacherLastName, &grade.SubjectName,
			&grade.SemesterCode, &grade.Signature, &isEdited, &isDeleted,
		); err != nil {
			return nil, err
		}

		// Convert int to bool
		grade.IsEdited = isEdited == 1
		grade.IsDeleted = isDeleted == 1

		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	pupilGradeGroup := tenantmodels.GradePupilGroup{
		Pupil:        *pupil,
		Grades:       grades,
		AverageGrade: CalculateAverageGrade(grades),
	}

	return &pupilGradeGroup, nil
}

// CreateGrade inserts a new grade into the database and returns the created grade.
func CreateGrade(
	grade *tenantmodels.Grade,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) (*tenantmodels.GradePupilGroup, error) {
	var err error

	if grade.Grade > 5 || grade.Grade < 1 {
		return nil, fmt.Errorf("ocjena mora biti između 1 i 5")
	}

	teacherForSignature, err := GetTeacherByID(
		fmt.Sprintf("%d", grade.TeacherID), workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `INSERT INTO student_grades (pupil_id, section_id, subject_code,
	grade, grade_date, type, teacher_id, semester_code, signature) VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?)`

	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	_, err = tx.Exec(query, grade.PupilID, grade.SectionID, grade.SubjectCode,
		grade.Grade, grade.GradeDate, grade.Type, grade.TeacherID, grade.SemesterCode,
		signature,
	)
	if err != nil {
		return nil, err
	}

	createdGrade, err := GetGradesForSectionSubjectPupil(
		grade.SectionID, grade.SubjectCode, grade.PupilID, grade.SemesterCode, tx, workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return createdGrade, nil
}

// DeleteGrade removes a grade from the database and returns the updated grades for the pupil.
func DeleteGrade(
	grade *tenantmodels.Grade,
	signature string,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) (*tenantmodels.GradePupilGroup, error) {
	var err error

	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	updateSignatureQuery := `UPDATE student_grades SET signature = ? WHERE id = ?`
	_, err = tx.Exec(updateSignatureQuery, signature, grade.ID)
	if err != nil {
		return nil, err
	}

	query := `DELETE FROM student_grades WHERE id = ?`
	_, err = tx.Exec(query, grade.ID)
	if err != nil {
		return nil, err
	}

	gradesAfterDeletion, err := GetGradesForSectionSubjectPupil(
		grade.SectionID, grade.SubjectCode, grade.PupilID, grade.SemesterCode, tx, workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return gradesAfterDeletion, nil
}

// EditGrade updates a grade
func EditGrade(
	grade *tenantmodels.Grade,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) (*tenantmodels.GradePupilGroup, error) {
	var err error

	if grade.Grade > 5 || grade.Grade < 1 {
		return nil, fmt.Errorf("ocjena mora biti između 1 i 5")
	}

	teacherForSignature, err := GetTeacherByID(
		fmt.Sprintf("%d", grade.TeacherID), workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	// Start a transaction
	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `UPDATE student_grades
	SET grade = ?, grade_date = ?, type = ?, teacher_id = ?, signature = ? WHERE id = ?`

	signature := teacherForSignature.Name + " " + teacherForSignature.LastName

	_, err = tx.Exec(query, grade.Grade, grade.GradeDate, grade.Type,
		grade.TeacherID, signature, grade.ID)
	if err != nil {
		return nil, err
	}

	updatedGradeData, err := GetGradesForSectionSubjectPupil(
		grade.SectionID, grade.SubjectCode, grade.PupilID, grade.SemesterCode, tx, workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updatedGradeData, nil
}

// GetSectionBehaviourGradesForPupilHelper retrieves all behaviour grades for a specific pupil
// within a given section, with one grade per semester that the section spans.
//
// Parameters:
//   - pupilID: unique identifier of the pupil
//   - sectionID: unique identifier of the section
//
// Returns:
//   - []BehaviourGrade: slice of behaviour grades, one per semester
//   - error: nil on success, error details on failure
//
// The returned slice length should equal the number of semesters in the section.
// Grades are ordered ascending by semester.
func GetSectionBehaviourGradesForPupilHelper(
	pupilID,
	sectionID int,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.BehaviourGrade, error) {
	query := `SELECT pb.id, pb.pupil_id, pb.section_id, pb.behaviour,
	pb.semester_code, s.semester_name, p.name, p.last_name,
	pb.signature
	FROM pupil_behaviour pb
	JOIN pupils p ON p.id = pb.pupil_id
	JOIN ednevnik_workspace.semester s ON s.semester_code = pb.semester_code
	WHERE pb.pupil_id = ? AND pb.section_id = ?
	ORDER BY s.progress_level ASC`

	rows, err := tenantDB.Query(query, pupilID, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var behaviourGrades []tenantmodels.BehaviourGrade
	for rows.Next() {
		var behaviourGrade tenantmodels.BehaviourGrade
		if err := rows.Scan(
			&behaviourGrade.ID, &behaviourGrade.PupilID, &behaviourGrade.SectionID,
			&behaviourGrade.Behaviour, &behaviourGrade.SemesterCode,
			&behaviourGrade.SemesterName, &behaviourGrade.Name, &behaviourGrade.LastName,
			&behaviourGrade.Signature,
		); err != nil {
			return nil, err
		}
		behaviourGrades = append(behaviourGrades, behaviourGrade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return behaviourGrades, nil
}

// GetPupilBehaviourGradeByID returns a behaviour grade by ID
func GetPupilBehaviourGradeByID(
	behaviourGradeID int,
	tenantDB *sql.DB,
) (*tenantmodels.BehaviourGrade, error) {
	query := `SELECT pb.id, pb.pupil_id, pb.section_id, pb.behaviour,
	pb.semester_code, s.semester_name, p.name, p.last_name,
	pb.signature
	FROM pupil_behaviour pb
	JOIN pupils p ON p.id = pb.pupil_id
	JOIN ednevnik_workspace.semester s ON s.semester_code = pb.semester_code
	WHERE pb.id = ?
	ORDER BY s.progress_level ASC`

	var behaviourGrade tenantmodels.BehaviourGrade

	err := tenantDB.QueryRow(query, behaviourGradeID).Scan(
		&behaviourGrade.ID, &behaviourGrade.PupilID, &behaviourGrade.SectionID,
		&behaviourGrade.Behaviour, &behaviourGrade.SemesterCode,
		&behaviourGrade.SemesterName, &behaviourGrade.Name, &behaviourGrade.LastName,
		&behaviourGrade.Signature,
	)
	if err != nil {
		return nil, err
	}

	return &behaviourGrade, nil
}

// GetPupilBehaviourGradeHistoryHelper returns history for a specific behaviour grade
func GetPupilBehaviourGradeHistoryHelper(
	behaviourGradeID int,
	tenantDB *sql.DB,
) ([]tenantmodels.BehaviourGrade, error) {
	query := `SELECT his.id, his.pupil_id, his.section_id, his.behaviour,
	his.semester_code, s.semester_name, p.name, p.last_name,
	his.signature, his.ROW_END
	FROM pupil_behaviour FOR SYSTEM_TIME ALL his
	JOIN pupils p ON p.id = his.pupil_id
	JOIN ednevnik_workspace.semester s ON s.semester_code = his.semester_code
	WHERE his.id = ? AND YEAR(his.ROW_END) < 2038
	ORDER BY his.ROW_END ASC`

	var behaviourGrades []tenantmodels.BehaviourGrade

	rows, err := tenantDB.Query(query, behaviourGradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var behaviourGrade tenantmodels.BehaviourGrade
		if err := rows.Scan(
			&behaviourGrade.ID, &behaviourGrade.PupilID, &behaviourGrade.SectionID,
			&behaviourGrade.Behaviour, &behaviourGrade.SemesterCode,
			&behaviourGrade.SemesterName, &behaviourGrade.Name, &behaviourGrade.LastName,
			&behaviourGrade.Signature, &behaviourGrade.ValidUntil,
		); err != nil {
			return nil, err
		}
		behaviourGrades = append(behaviourGrades, behaviourGrade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return behaviourGrades, nil
}

// GetAllFinalGradesForSection all final grades for section
func GetAllFinalGradesForSection(
	sectionID int, tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Grade, error) {
	query := `SELECT sg.id, sg.pupil_id, sg.section_id, sg.subject_code,
	sg.grade, sg.grade_date, sg.type, sg.semester_code, sg.signature
	FROM student_grades sg
	WHERE sg.section_id = ? AND sg.type = 'final'`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var grades []tenantmodels.Grade
	for rows.Next() {
		var grade tenantmodels.Grade
		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.Type,
			&grade.SemesterCode, &grade.Signature,
		); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

// GetGradeEditHistoryHelper retrieves the edit history for a specific grade ID,
// excluding the current version. Returns historical versions ordered by
// modification time (oldest first).
func GetGradeEditHistoryHelper(
	gradeID int,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Grade, error) {
	query := `SELECT hist.id, hist.pupil_id, hist.section_id, hist.subject_code,
	hist.grade, hist.grade_date, hist.type, COALESCE(t.name, ''), COALESCE(t.last_name, ''),
	s.subject_name, hist.semester_code, hist.signature,
	1 AS is_edited, hist.ROW_END
	FROM student_grades FOR SYSTEM_TIME ALL hist
	JOIN ednevnik_workspace.subjects s ON hist.subject_code = s.subject_code
	LEFT JOIN ednevnik_workspace.teachers t ON hist.teacher_id = t.id
	WHERE hist.id = ?
	AND YEAR(hist.ROW_END) < 2038
	ORDER BY hist.ROW_START ASC`

	rows, err := tenantDB.Query(query, gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []tenantmodels.Grade
	for rows.Next() {
		var grade tenantmodels.Grade

		if err := rows.Scan(
			&grade.ID, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.Type,
			&grade.TeacherName, &grade.TeacherLastName, &grade.SubjectName,
			&grade.SemesterCode, &grade.Signature, &grade.IsEdited,
			&grade.ValidUntil,
		); err != nil {
			return nil, err
		}
		grades = append(grades, grade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

// GetPupilGradesForCompleteGradebook returns all grades for a gradebook
func GetGradeDataForCompleteGradebook(
	sectionID int,
	semesters []wpmodels.TenantSemester,
	pupils []tenantmodels.Pupil,
	subjects []wpmodels.Subject,
	tenantDB *sql.DB,
) ([]tenantmodels.CompleteGradebookData, error) {
	// Get all grades for section
	gradesQuery := `SELECT DISTINCT sg.id, sg.type, sg.pupil_id, sg.section_id,
    sg.subject_code, sg.grade, sg.grade_date, subj.subject_name, sg.signature,
	sg.semester_code
	FROM student_grades sg
	JOIN ednevnik_workspace.subjects subj ON subj.subject_code = sg.subject_code
	WHERE sg.section_id = ?
	ORDER BY sg.grade_date ASC, subj.subject_name ASC`
	gradeRows, err := tenantDB.Query(gradesQuery, sectionID)
	if err != nil {
		return nil, err
	}
	defer gradeRows.Close()

	// Create lookup map for subject grades
	gradesByPupilSemesterSubject := make(map[string][]tenantmodels.Grade)

	for gradeRows.Next() {
		var grade tenantmodels.Grade
		if err := gradeRows.Scan(
			&grade.ID, &grade.Type, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate, &grade.SubjectName,
			&grade.Signature, &grade.SemesterCode,
		); err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%d-%s-%s", grade.PupilID, grade.SemesterCode, grade.SubjectCode)
		gradesByPupilSemesterSubject[key] = append(gradesByPupilSemesterSubject[key], grade)
	}
	if err := gradeRows.Err(); err != nil {
		return nil, err
	}

	// Get all behaviour grades for sections
	behaviourQuery := `SELECT b.id, b.pupil_id, b.section_id, b.behaviour,
	b.semester_code, b.row_start, b.signature FROM pupil_behaviour FOR SYSTEM_TIME ALL b
	WHERE b.section_id = ? ORDER BY b.row_start ASC`
	behaviourRows, err := tenantDB.Query(behaviourQuery, sectionID)
	if err != nil {
		return nil, err
	}
	defer behaviourRows.Close()

	// Create lookup map for behavior grades
	behaviourByPupilSemester := make(map[string][]tenantmodels.BehaviourGrade)

	for behaviourRows.Next() {
		var behaviourGrade tenantmodels.BehaviourGrade
		if err := behaviourRows.Scan(
			&behaviourGrade.ID, &behaviourGrade.PupilID, &behaviourGrade.SectionID,
			&behaviourGrade.Behaviour, &behaviourGrade.SemesterCode,
			&behaviourGrade.Date, &behaviourGrade.Signature,
		); err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%d-%s", behaviourGrade.PupilID, behaviourGrade.SemesterCode)
		behaviourByPupilSemester[key] = append(behaviourByPupilSemester[key], behaviourGrade)
	}
	if err := behaviourRows.Err(); err != nil {
		return nil, err
	}

	completeGradebookData := make([]tenantmodels.CompleteGradebookData, 0, len(pupils))

	for _, pupil := range pupils {
		var completeGradebookDataForPupil tenantmodels.CompleteGradebookData
		completeGradebookDataForPupil.PupilName = pupil.LastName + " " + pupil.Name
		completeGradebookDataForPupil.PupilUnenrolled = pupil.Unenrolled

		completeGradebookDataForPupil.GradesForSemester = make([]tenantmodels.SemesterGradeGroup, 0, len(semesters))

		for _, semester := range semesters {
			var semesterGradesForPupil tenantmodels.SemesterGradeGroup
			semesterGradesForPupil.SemesterName = semester.SemesterName

			semesterGradesForPupil.SubjectGrades = make([]tenantmodels.SubjectGradeGroup, 0, len(subjects))

			for _, subject := range subjects {
				var subjectGradesForPupil tenantmodels.SubjectGradeGroup
				subjectGradesForPupil.SubjectName = subject.SubjectName
				subjectGradesForPupil.SubjectCode = subject.SubjectCode

				gradeKey := fmt.Sprintf(
					"%d-%s-%s", pupil.ID, semester.SemesterCode, subject.SubjectCode,
				)
				if grades, exists := gradesByPupilSemesterSubject[gradeKey]; exists {
					subjectGradesForPupil.Grades = grades
				}

				semesterGradesForPupil.SubjectGrades = append(
					semesterGradesForPupil.SubjectGrades, subjectGradesForPupil,
				)
			}

			behaviourKey := fmt.Sprintf("%d-%s", pupil.ID, semester.SemesterCode)
			if behaviourGrades, exists := behaviourByPupilSemester[behaviourKey]; exists {
				semesterGradesForPupil.BehaviourGrades = behaviourGrades
			}

			completeGradebookDataForPupil.GradesForSemester = append(
				completeGradebookDataForPupil.GradesForSemester, semesterGradesForPupil,
			)
		}
		completeGradebookData = append(completeGradebookData, completeGradebookDataForPupil)
	}

	return completeGradebookData, nil
}
