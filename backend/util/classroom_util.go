package util

import (
	"database/sql"
	tenantmodels "ednevnik-backend/models/tenant"
)

// CreateClassroom TODO: Add description
func CreateClassroom(
	data tenantmodels.Classroom, tenantDB *sql.DB,
) error {
	query := `INSERT INTO classroom (code, capacity, type) VALUES (?, ?, ?)`
	_, err := tenantDB.Exec(query, data.Code, data.Capacity, data.Type)
	return err
}

// UpdateClassroom TODO: Add description
func UpdateClassroom(
	data tenantmodels.Classroom, tenantDB *sql.DB, oldCode string,
) error {
	query := `UPDATE classroom SET capacity = ?, code = ?, type = ? WHERE code = ?`
	_, err := tenantDB.Exec(query, data.Capacity, data.Code, data.Type, oldCode)
	return err
}

// GetAllClassroomsForTenant TODO: Add description
func GetAllClassroomsForTenant(
	tenantDB *sql.DB,
) ([]tenantmodels.Classroom, error) {
	query := `SELECT code, capacity, type, CONCAT(type, ' ', code) as name FROM classroom`
	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classrooms []tenantmodels.Classroom
	for rows.Next() {
		var classroom tenantmodels.Classroom
		if err := rows.Scan(&classroom.Code, &classroom.Capacity, &classroom.Type,
			&classroom.Name); err != nil {
			return nil, err
		}
		classrooms = append(classrooms, classroom)
	}
	return classrooms, nil
}

// DeleteClassroom TODO: Add description
func DeleteClassroom(
	code string, tenantDB *sql.DB,
) error {
	query := `DELETE FROM classroom WHERE code = ?`
	_, err := tenantDB.Exec(query, code)
	return err
}
