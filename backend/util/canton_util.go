package util

import (
	"database/sql"
	wpmodels "ednevnik-backend/models/workspace"
)

// GetAllCantons fetches all cantons from the workspace database
func GetAllCantons(db *sql.DB) ([]wpmodels.Canton, error) {
	rows, err := db.Query("SELECT canton_code, canton_name, country FROM cantons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cantons []wpmodels.Canton
	for rows.Next() {
		var c wpmodels.Canton
		if err := rows.Scan(&c.CantonCode, &c.CantonName, &c.Country); err != nil {
			return nil, err
		}
		cantons = append(cantons, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cantons, nil
}
