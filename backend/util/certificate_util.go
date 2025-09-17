package util

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	"fmt"
	"math"
)

// GetCertificateData retrieves the certificate data for a pupil in a section
// It includes final grades and behaviour grades for the semester with the
// highest progress level.
func GetCertificateData(
	tenantID,
	sectionID,
	pupilID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) (*commonmodels.Certificate, error) {
	tenant, err := GetTenantByID(
		fmt.Sprintf("%d", tenantID), workspaceDB,
	)
	if err != nil {
		return nil, err
	}

	section, err := GetSectionByID(int64(sectionID), tenantDB)
	if err != nil {
		return nil, err
	}

	pupil, err := GetTenantPupilByID(pupilID, tenantDB)
	if err != nil {
		return nil, err
	}

	finalGradesQuery := `SELECT sg.id, sg.type, sg.pupil_id, sg.section_id,
    sg.subject_code, sg.grade, sg.grade_date, subj.subject_name, sg.signature
	FROM student_grades sg
	JOIN ednevnik_workspace.semester sem ON sg.semester_code = sem.semester_code
	JOIN ednevnik_workspace.subjects subj ON subj.subject_code = sg.subject_code
	WHERE sg.pupil_id = ?  AND sg.section_id = ? AND sg.type = 'final'
    AND sem.progress_level = (SELECT MAX(progress_level) FROM ednevnik_workspace.semester)
	ORDER BY subj.subject_name ASC`

	finalGradeRows, err := tenantDB.Query(finalGradesQuery, pupilID, sectionID)
	if err != nil {
		return nil, err
	}
	defer finalGradeRows.Close()

	var finalGrades []tenantmodels.Grade
	for finalGradeRows.Next() {
		var grade tenantmodels.Grade
		err := finalGradeRows.Scan(
			&grade.ID, &grade.Type, &grade.PupilID, &grade.SectionID,
			&grade.SubjectCode, &grade.Grade, &grade.GradeDate,
			&grade.SubjectName, &grade.Signature,
		)
		if err != nil {
			return nil, err
		}
		finalGrades = append(finalGrades, grade)
	}

	averageFinalGrade := CalculateAverageFinalGrade(finalGrades)

	behaviourQuery := `SELECT b.id, b.pupil_id, b.section_id, b.behaviour,
	b.semester_code FROM pupil_behaviour b
	JOIN ednevnik_workspace.semester sem ON b.semester_code = sem.semester_code
	WHERE b.pupil_id = ? AND b.section_id = ?
	AND sem.progress_level = (SELECT MAX(progress_level) FROM ednevnik_workspace.semester)`

	var behaviour tenantmodels.BehaviourGrade
	err = tenantDB.QueryRow(behaviourQuery, pupilID, sectionID).Scan(
		&behaviour.ID, &behaviour.PupilID, &behaviour.SectionID,
		&behaviour.Behaviour, &behaviour.SemesterCode,
	)
	if err != nil {
		return nil, err
	}

	// Graduate grade round averageFinalGrade to 2 decimals
	graduateGrade := int(math.Round(averageFinalGrade))

	passed := true
	for _, finalGrade := range finalGrades {
		if finalGrade.Grade < 2 {
			passed = false
			break
		}
	}

	certificate := &commonmodels.Certificate{
		Tenant:         *tenant,
		Section:        section,
		Pupil:          *pupil,
		FinalGrades:    finalGrades,
		BehaviourGrade: behaviour,
		AverageGrade:   averageFinalGrade,
		GraduateGrade:  graduateGrade,
		Passed:         passed,
	}

	return certificate, nil
}
