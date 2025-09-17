package util

import (
	"database/sql"
	"ednevnik-backend/models/interfaces"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"strconv"
)

// GetAllNPPSemesters TODO: Add description
func GetAllNPPSemesters(workspaceDB *sql.DB) ([]wpmodels.NPPSemester, error) {
	query := `SELECT npp.npp_code, npp.npp_name, ns.semester_code, s.semester_name,
	ns.start_date, ns.end_date
    FROM npp_semester ns
    JOIN npp ON ns.npp_code = npp.npp_code
    JOIN semester s ON ns.semester_code = s.semester_code
    ORDER BY npp.npp_code, s.progress_level`

	rows, err := workspaceDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nppSemesters []wpmodels.NPPSemester

	for rows.Next() {
		var nppSemester wpmodels.NPPSemester
		err := rows.Scan(
			&nppSemester.NPPCode,
			&nppSemester.NPPName,
			&nppSemester.SemesterCode,
			&nppSemester.SemesterName,
			&nppSemester.StartDate,
			&nppSemester.EndDate,
		)
		if err != nil {
			return nil, err
		}
		nppSemester.FullName = nppSemester.GetFullName()
		nppSemesters = append(nppSemesters, nppSemester)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nppSemesters, nil
}

// UpdateNPPSemesterDates TODO: Add description
func UpdateNPPSemesterDates(
	workspaceDB *sql.DB, nppCode string, semesterCode string, startDate, endDate string,
) (wpmodels.NPPSemester, error) {
	query := `UPDATE npp_semester SET start_date = ?, end_date = ?
    WHERE npp_code = ? AND semester_code = ?`

	_, err := workspaceDB.Exec(query, startDate, endDate, nppCode, semesterCode)
	if err != nil {
		return wpmodels.NPPSemester{}, err
	}

	updatedSemester, err := GetNPPSemesterByCodes(
		workspaceDB, nppCode, semesterCode,
	)
	if err != nil {
		return wpmodels.NPPSemester{}, err
	}

	return *updatedSemester, nil
}

// GetNPPSemesterByCodes TODO: Add description
func GetNPPSemesterByCodes(
	workspaceDB *sql.DB, nppCode string, semesterCode string,
) (*wpmodels.NPPSemester, error) {
	query := `SELECT npp.npp_code, npp.npp_name, ns.semester_code, s.semester_name,
    ns.start_date, ns.end_date
    FROM npp_semester ns
    JOIN npp ON ns.npp_code = npp.npp_code
    JOIN semester s ON ns.semester_code = s.semester_code
    WHERE npp.npp_code = ? AND ns.semester_code = ?`

	row := workspaceDB.QueryRow(query, nppCode, semesterCode)

	var nppSemester wpmodels.NPPSemester
	err := row.Scan(
		&nppSemester.NPPCode,
		&nppSemester.NPPName,
		&nppSemester.SemesterCode,
		&nppSemester.SemesterName,
		&nppSemester.StartDate,
		&nppSemester.EndDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	nppSemester.FullName = nppSemester.GetFullName()

	return &nppSemester, nil
}

// TenantSemesterCleanup TODO: Add description
func TenantSemesterCleanup(
	tenantID string, workspaceDB *sql.DB,
) error {
	// Delete tenant semester records who do not have a corresponding
	// tenant curriculum record
	query := `
        DELETE FROM tenant_semester
        WHERE tenant_id = ?
        AND npp_code NOT IN (
            SELECT c.npp_code
            FROM curriculum_tenant ct
            JOIN curriculum c ON ct.curriculum_code = c.curriculum_code
            WHERE ct.tenant_id = ?
        )`

	_, err := workspaceDB.Exec(query, tenantID, tenantID)
	if err != nil {
		return fmt.Errorf("error cleaning up tenant semster records: %v", err)
	}

	return nil
}

// TenantSemesterAssign TODO: Add description
func TenantSemesterAssign(
	tenantID string, workspaceDB *sql.DB,
) error {
	// Insert missing tenant_semester records for the given tenant.
	// For each curriculum assigned to the tenant, find all semesters from npp_semester
	// and insert them into tenant_semester if not already present.
	// Uses INSERT IGNORE to avoid duplicate (tenant_id, semester_code) entries.
	query := `
        INSERT IGNORE INTO tenant_semester (tenant_id, semester_code, start_date, end_date,
		npp_code)
        SELECT ct.tenant_id, ns.semester_code, ns.start_date, ns.end_date, ns.npp_code
        FROM curriculum_tenant ct
        JOIN curriculum c ON ct.curriculum_code = c.curriculum_code
        JOIN npp_semester ns ON ns.npp_code = c.npp_code
        WHERE ct.tenant_id = ?
    `
	_, err := workspaceDB.Exec(query, tenantID)
	if err != nil {
		return fmt.Errorf("error assigning tenant semester records: %v", err)
	}
	return nil
}

// GetSemestersForTenant TODO: Add description
func GetSemestersForTenant(workspaceDB *sql.DB, tenantID string) ([]wpmodels.TenantSemester, error) {
	query := `SELECT ts.tenant_id, npp.npp_name, ts.semester_code, s.semester_name,
    ts.start_date, ts.end_date, ts.npp_code
    FROM tenant_semester ts
    JOIN npp ON ts.npp_code = npp.npp_code
    JOIN semester s ON ts.semester_code = s.semester_code
    WHERE ts.tenant_id = ?
    ORDER BY ts.npp_code, s.progress_level`

	rows, err := workspaceDB.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenantSemesters []wpmodels.TenantSemester

	for rows.Next() {
		var tenantSemester wpmodels.TenantSemester
		err := rows.Scan(
			&tenantSemester.TenantID,
			&tenantSemester.NPPName,
			&tenantSemester.SemesterCode,
			&tenantSemester.SemesterName,
			&tenantSemester.StartDate,
			&tenantSemester.EndDate,
			&tenantSemester.NPPCode,
		)
		if err != nil {
			return nil, err
		}
		tenantSemester.FullName = tenantSemester.GetFullName()
		tenantSemesters = append(tenantSemesters, tenantSemester)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenantSemesters, nil
}

// GetSemestersForSectionHelper returns semesters (polugodišta) for a section
func GetSemestersForSectionHelper(
	workspaceDB *sql.DB,
	tenantDB interfaces.DatabaseQuerier,
	tenantID,
	sectionID string,
) ([]wpmodels.TenantSemester, error) {
	tenantSemesters, err := GetSemestersForTenant(workspaceDB, tenantID)
	if err != nil {
		return nil, err
	}

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		return nil, err
	}

	section, err := GetSectionByID(
		int64(sectionIDInt),
		tenantDB,
	)
	if err != nil {
		return nil, err
	}

	curriculum, err := GetCurriculumByCode(
		workspaceDB, section.CurriculumCode,
	)
	if err != nil {
		return nil, err
	}

	var sectionSemester []wpmodels.TenantSemester

	for _, tenantSemester := range tenantSemesters {
		if tenantSemester.NPPCode == curriculum.NPPCode {
			sectionSemester = append(sectionSemester, tenantSemester)
		}
	}

	return sectionSemester, nil
}

// UpdateTenantSemesterDates TODO: Add description
func UpdateTenantSemesterDates(
	workspaceDB *sql.DB, tenantID, semesterCode, startDate, endDate, nppCode string,
) (wpmodels.TenantSemester, error) {
	query := `UPDATE tenant_semester SET start_date = ?, end_date = ?
    WHERE tenant_id = ? AND semester_code = ? AND npp_code = ?`

	_, err := workspaceDB.Exec(query, startDate, endDate, tenantID, semesterCode, nppCode)
	if err != nil {
		return wpmodels.TenantSemester{}, err
	}

	updatedSemester, err := GetTenantSemesterByCodes(
		workspaceDB, tenantID, semesterCode, nppCode,
	)
	if err != nil {
		return wpmodels.TenantSemester{}, err
	}

	return *updatedSemester, nil
}

// GetTenantSemesterByCodes TODO: Add description
func GetTenantSemesterByCodes(
	workspaceDB *sql.DB, tenantID, semesterCode, nppCode string,
) (*wpmodels.TenantSemester, error) {
	query := `SELECT ts.tenant_id, npp.npp_name, ts.semester_code, s.semester_name,
    ts.start_date, ts.end_date, ts.npp_code
    FROM tenant_semester ts
    JOIN npp ON ts.npp_code = npp.npp_code
    JOIN semester s ON ts.semester_code = s.semester_code
    WHERE ts.tenant_id = ? AND ts.semester_code = ? AND ts.npp_code = ?`

	row := workspaceDB.QueryRow(query, tenantID, semesterCode, nppCode)

	var tenantSemester wpmodels.TenantSemester
	err := row.Scan(
		&tenantSemester.TenantID,
		&tenantSemester.NPPName,
		&tenantSemester.SemesterCode,
		&tenantSemester.SemesterName,
		&tenantSemester.StartDate,
		&tenantSemester.EndDate,
		&tenantSemester.NPPCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	tenantSemester.FullName = tenantSemester.GetFullName()

	return &tenantSemester, nil
}
