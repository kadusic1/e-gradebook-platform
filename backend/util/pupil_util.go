package util

import (
	"database/sql"
	"ednevnik-backend/models/interfaces"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ListPupilAccounts TODO: Add description
func ListPupilAccounts(
	workspaceDB *sql.DB, loggedInTeacher wpmodels.Claims,
) ([]tenantmodels.Pupil, error) {
	var query string
	var rows *sql.Rows
	var err error
	if loggedInTeacher.AccountType == "root" {
		query = `SELECT p.id, p.name, p.last_name, p.gender, p.address,
		p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
		p.religion, a.email, p.place_of_birth
		FROM pupil_global p
		JOIN accounts a ON p.account_id = a.id
		WHERE a.account_type = 'pupil'`
		rows, err = workspaceDB.Query(query)
	} else {
		query = `SELECT p.id, p.name, p.last_name, p.gender, p.address,
		p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
		p.religion, a.email, p.place_of_birth
		FROM pupil_global p
		JOIN accounts a ON p.account_id = a.id
		WHERE a.created_by_teacher_id = ?
		AND a.account_type = 'pupil'`
		rows, err = workspaceDB.Query(query, loggedInTeacher.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("error listing pupil accounts: %v", err)
	}

	defer rows.Close()
	pupils := []tenantmodels.Pupil{}
	for rows.Next() {
		var pupil tenantmodels.Pupil
		err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.LastName,
			&pupil.Gender,
			&pupil.Address,
			&pupil.GuardianName,
			&pupil.PhoneNumber,
			&pupil.GuardianNumber,
			&pupil.DateOfBirth,
			&pupil.Religion,
			&pupil.Email,
			&pupil.PlaceOfBirth,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning pupil row: %v", err)
		}
		pupils = append(pupils, pupil)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	return pupils, nil
}

// GetPupilsForSection TODO: Add description
func GetPupilsForSection(
	sectionID string,
	includeUnenrolled bool,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Pupil, error) {
	pupils := []tenantmodels.Pupil{}

	var query string
	if includeUnenrolled {
		query = `SELECT p.id, p.name, p.last_name, p.gender, p.address,
		p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
		p.religion, a.email, p.place_of_birth, ps.is_active
		FROM pupils p
		JOIN pupils_sections ps ON p.id = ps.pupil_id
		JOIN sections s ON ps.section_id = s.id
		JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
		WHERE ps.section_id = ? ORDER BY ps.is_active DESC, p.last_name ASC,
		p.name ASC`
	} else {
		query = `SELECT p.id, p.name, p.last_name, p.gender, p.address,
		p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
		p.religion, a.email, p.place_of_birth, ps.is_active
		FROM pupils p
		JOIN pupils_sections ps ON p.id = ps.pupil_id
		JOIN sections s ON ps.section_id = s.id
		JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
		WHERE ps.section_id = ? AND ps.is_active = 1
		ORDER BY p.last_name ASC, p.name ASC`
	}

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pupil tenantmodels.Pupil
		var isActive bool
		if err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.LastName,
			&pupil.Gender,
			&pupil.Address,
			&pupil.GuardianName,
			&pupil.PhoneNumber,
			&pupil.GuardianNumber,
			&pupil.DateOfBirth,
			&pupil.Religion,
			&pupil.Email,
			&pupil.PlaceOfBirth,
			&isActive,
		); err != nil {
			return nil, err
		}
		pupil.Unenrolled = !isActive
		pupils = append(pupils, pupil)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pupils, nil
}

// GetPupilsForCompleteGradebook returns pupil data necessary for complete gradebook
func GetPupilsForCompleteGradebook(
	sectionID string,
	tenantDB interfaces.DatabaseQuerier,
) ([]tenantmodels.Pupil, error) {
	pupils := []tenantmodels.Pupil{}

	query := `SELECT p.id, p.name, p.last_name, p.guardian_name, p.phone_number,
	a.email, p.religion, p.is_commuter, ps.is_active
	FROM ednevnik_workspace.pupil_global p
	JOIN pupils_sections ps ON p.id = ps.pupil_id
	JOIN sections s ON ps.section_id = s.id
	JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
	WHERE ps.section_id = ?
	ORDER BY ps.is_active DESC, p.last_name ASC, p.name ASC`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var isCommuter interface{}
		var pupil tenantmodels.Pupil
		var isActive bool
		if err := rows.Scan(
			&pupil.ID,
			&pupil.Name,
			&pupil.LastName,
			&pupil.GuardianName,
			&pupil.PhoneNumber,
			&pupil.Email,
			&pupil.Religion,
			&isCommuter,
			&isActive,
		); err != nil {
			return nil, err
		}
		switch v := isCommuter.(type) {
		case nil:
			pupil.IsCommuter = "Nije postavljeno"
		case int64:
			if v == 1 {
				pupil.IsCommuter = "Da"
			} else {
				pupil.IsCommuter = "Ne"
			}
		default:
			pupil.IsCommuter = "Ne"
		}
		pupil.Unenrolled = !isActive
		pupils = append(pupils, pupil)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pupils, nil
}

// CreatePupil TODO: Add description
func CreatePupil(
	pupil tenantmodels.Pupil,
	workspaceDB *sql.DB,
) (tenantmodels.Pupil, error) {
	var err error
	hash, err := bcrypt.GenerateFromPassword([]byte(pupil.Password), bcrypt.DefaultCost)
	if err != nil {
		return tenantmodels.Pupil{}, err
	}
	pupil.Password = string(hash)

	tx, err := workspaceDB.Begin()
	if err != nil {
		return tenantmodels.Pupil{}, fmt.Errorf("error starting transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// First insert account
	accountQuery := `INSERT INTO accounts (email, password, account_type)
	VALUES (?, ?, 'pupil')`

	res, err := tx.Exec(
		accountQuery,
		pupil.Email,
		pupil.Password,
	)
	if err != nil {
		if IsDuplicateEmailError(err) {
			return tenantmodels.Pupil{}, fmt.Errorf("korisnik sa ovim emailom već postoji")
		}
		return tenantmodels.Pupil{}, err
	}

	accountID, err := res.LastInsertId()
	if err != nil {
		return tenantmodels.Pupil{}, err
	}

	// Insert global pupil into workspace DB
	globalPupilQuery := `INSERT INTO pupil_global (name, last_name, gender, address,
	guardian_name, phone_number, guardian_number, date_of_birth, religion,
	account_id, place_of_birth
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err = tx.Exec(
		globalPupilQuery,
		pupil.Name,
		pupil.LastName,
		pupil.Gender,
		pupil.Address,
		pupil.GuardianName,
		pupil.PhoneNumber,
		pupil.GuardianNumber,
		pupil.DateOfBirth,
		pupil.Religion,
		accountID,
		pupil.PlaceOfBirth,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return tenantmodels.Pupil{}, fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
		}
		return tenantmodels.Pupil{}, err
	}

	globalPupilID, err := res.LastInsertId()
	if err != nil {
		return tenantmodels.Pupil{}, err
	}
	pupil.ID = int(globalPupilID)
	pupil.Password = "" // Clear password after hashing

	err = tx.Commit()
	if err != nil {
		return tenantmodels.Pupil{}, fmt.Errorf("error committing transaction: %v", err)
	}

	return pupil, nil
}

// RegisterPupil TODO: Add description
func RegisterPupil(
	pupil tenantmodels.Pupil,
	workspaceDB *sql.DB,
) error {
	var err error
	domains, err := GetAllDomainsHelper(workspaceDB)
	if err != nil {
		return err
	}

	// Check if the email ends with a valid domain
	validDomain := false
	for _, domain := range domains {
		if strings.HasSuffix(pupil.Email, domain.Domain) {
			validDomain = true
			break
		}
	}
	if !validDomain {
		return fmt.Errorf("email učenika mora biti iz validne domene")
	}

	exists, err := AccountWithEmailExists(pupil.Email, workspaceDB)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("korisnik sa ovim emailom već postoji")
	}

	exists, err = PupilWithPhoneExists(pupil.PhoneNumber, workspaceDB)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
	}

	if err := ValidateIdentifier(pupil.Email); err != nil {
		return fmt.Errorf("ovaj email nije validan")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pupil.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	pupil.Password = string(hash)

	tx, err := workspaceDB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// First insert account
	accountQuery := `INSERT INTO pending_accounts (email, password, account_type)
	VALUES (?, ?, 'pupil')`

	res, err := tx.Exec(
		accountQuery,
		pupil.Email,
		pupil.Password,
	)
	if err != nil {
		if IsDuplicateEmailError(err) {
			return fmt.Errorf(
				"neverifikovani korisnik sa ovim emailom već postoji - provjerite email za aktivaciju",
			)
		}
		return err
	}

	accountID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Insert global pupil into workspace DB
	globalPupilQuery := `INSERT INTO pending_pupil_global (name, last_name, gender, address,
	guardian_name, phone_number, guardian_number, date_of_birth, religion,
	account_id, place_of_birth
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(
		globalPupilQuery,
		pupil.Name,
		pupil.LastName,
		pupil.Gender,
		pupil.Address,
		pupil.GuardianName,
		pupil.PhoneNumber,
		pupil.GuardianNumber,
		pupil.DateOfBirth,
		pupil.Religion,
		accountID,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return fmt.Errorf(
				"neverifikovani korisnik sa ovim brojem telefona već postoji - provjerite email za aktivaciju",
			)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	token, err := GetPendingAccountVerificationToken(int(accountID), workspaceDB)
	if err != nil {
		return fmt.Errorf("error getting verification token: %v", err)
	}

	go func() {
		_ = SendVerificationEmail(
			pupil.Email,
			fmt.Sprintf("%s %s", pupil.Name, pupil.LastName),
			fmt.Sprintf("%s/verify?token=%s", os.Getenv("FRONTEND_URL"), token),
		)
	}()

	return nil
}

// GetGlobalPupilByEmail TODO: Add description
func GetGlobalPupilByEmail(
	email string, workspaceDB *sql.DB,
) (*tenantmodels.Pupil, error) {
	var pupil tenantmodels.Pupil

	query := `SELECT pg.id, pg.name, pg.last_name, pg.gender, pg.address,
	pg.guardian_name, pg.phone_number, pg.guardian_number, pg.date_of_birth,
	pg.religion, a.password, a.email, pg.place_of_birth
	FROM pupil_global pg
	JOIN accounts a ON pg.account_id = a.id
	WHERE email = ?`

	err := workspaceDB.QueryRow(query, email).Scan(
		&pupil.ID,
		&pupil.Name,
		&pupil.LastName,
		&pupil.Gender,
		&pupil.Address,
		&pupil.GuardianName,
		&pupil.PhoneNumber,
		&pupil.GuardianNumber,
		&pupil.DateOfBirth,
		&pupil.Religion,
		&pupil.Password,
		&pupil.Email,
		&pupil.PlaceOfBirth,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pupil, nil
}

// GetGlobalPupilByEmail TODO: Add description
func GetGlobalPupilByParentAccessCode(
	parentAccessCode string, workspaceDB *sql.DB,
) (*tenantmodels.Pupil, error) {
	var pupil tenantmodels.Pupil

	query := `SELECT pg.id, pg.name, pg.last_name, pg.gender, pg.address,
	pg.guardian_name, pg.phone_number, pg.guardian_number, pg.date_of_birth,
	pg.religion, a.password, a.email, pg.place_of_birth
	FROM pupil_global pg
	JOIN accounts a ON pg.account_id = a.id
	WHERE pg.parent_access_code = ?`

	err := workspaceDB.QueryRow(query, parentAccessCode).Scan(
		&pupil.ID,
		&pupil.Name,
		&pupil.LastName,
		&pupil.Gender,
		&pupil.Address,
		&pupil.GuardianName,
		&pupil.PhoneNumber,
		&pupil.GuardianNumber,
		&pupil.DateOfBirth,
		&pupil.Religion,
		&pupil.Password,
		&pupil.Email,
		&pupil.PlaceOfBirth,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pupil, nil
}

// GetGlobalPupilByID TODO: Add description
func GetGlobalPupilByID(
	id string, workspaceDB interfaces.DatabaseQuerier,
) (*tenantmodels.Pupil, error) {
	var pupil tenantmodels.Pupil

	query := `SELECT p.id, p.name, p.last_name, p.gender, p.address,
	p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
	p.religion, a.password, a.email, p.place_of_birth, p.parent_access_code
	FROM pupil_global p
	JOIN accounts a ON p.account_id = a.id
	WHERE p.id = ?`

	err := workspaceDB.QueryRow(query, id).Scan(
		&pupil.ID,
		&pupil.Name,
		&pupil.LastName,
		&pupil.Gender,
		&pupil.Address,
		&pupil.GuardianName,
		&pupil.PhoneNumber,
		&pupil.GuardianNumber,
		&pupil.DateOfBirth,
		&pupil.Religion,
		&pupil.Password,
		&pupil.Email,
		&pupil.PlaceOfBirth,
		&pupil.ParentAccessCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pupil, nil
}

// GetTenantPupilByID finds a pupil by ID in the tenant database
func GetTenantPupilByID(
	pupilID int, tenantDB *sql.DB,
) (*tenantmodels.Pupil, error) {
	var pupil tenantmodels.Pupil

	query := `SELECT p.id, p.name, p.last_name, p.gender, p.address,
	p.guardian_name, p.phone_number, p.guardian_number, p.date_of_birth,
	p.religion, a.email, p.place_of_birth
	FROM pupils p
	JOIN ednevnik_workspace.accounts a ON p.account_id = a.id
	WHERE p.id = ?`

	err := tenantDB.QueryRow(query, pupilID).Scan(
		&pupil.ID,
		&pupil.Name,
		&pupil.LastName,
		&pupil.Gender,
		&pupil.Address,
		&pupil.GuardianName,
		&pupil.PhoneNumber,
		&pupil.GuardianNumber,
		&pupil.DateOfBirth,
		&pupil.Religion,
		&pupil.Email,
		&pupil.PlaceOfBirth,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &pupil, nil
}

// GetPupilTablePrivileges TODO: Add description
func GetPupilTablePrivileges() []struct {
	Table   string
	Actions string
} {
	return []struct {
		Table   string
		Actions string
	}{
		{"pupils", "SELECT"},
		{"sections", "SELECT"},
		{"pupils_sections", "SELECT"},
		{"student_grades", "SELECT"},
		{"pupils_sections_invite", "SELECT"},
		{"homeroom_assignments", "SELECT"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT"},
		{"pupil_attendance", "SELECT"},
		{"pupil_behaviour", "SELECT"},
	}
}

// DeleteGlobalPupilRecord TODO: Add description
func DeleteGlobalPupilRecord(pupil tenantmodels.Pupil, workspaceDB *sql.DB) error {
	pupilAccountID, err := pupil.GetAccountID(workspaceDB)
	if err != nil {
		return fmt.Errorf("error getting pupil account ID: %v", err)
	}

	deleteQuery := `DELETE FROM accounts WHERE id = ?;`

	_, err = workspaceDB.Exec(deleteQuery, pupilAccountID)
	if err != nil {
		return fmt.Errorf("error deleting global pupil record: %v", err)
	}

	return nil
}

// DeleteTenantPupilRecord TODO: Add description
func DeleteTenantPupilRecord(
	pupilID string, tenantID int, workspaceDB *sql.DB, tenantDB *sql.DB,
) error {
	tenantPupilRecordDeleteQuery := `DELETE FROM pupil_tenant WHERE
	pupil_id = ? AND tenant_id = ?`
	_, err := workspaceDB.Exec(tenantPupilRecordDeleteQuery, pupilID, tenantID)
	if err != nil {
		return fmt.Errorf("error deleting tenant pupil record: %v", err)
	}

	tenantDbPupilDeleteQuery := "DELETE FROM pupils WHERE id = ?"
	_, err = tenantDB.Exec(tenantDbPupilDeleteQuery, pupilID)
	if err != nil {
		return fmt.Errorf("error deleting pupil from tenant db: %v", err)
	}

	// No need to delete section invites we want to keep them for history
	// tenant_sections_invite_deleteQuery := `DELETE FROM pupils_sections_invite
	// WHERE pupil_id = ?`
	// _, err = tenantDB.Exec(tenant_sections_invite_deleteQuery, pupil_id)
	// if err != nil {
	// 	return fmt.Errorf("error deleting pupil section invites from tenant db: %v", err)
	// }

	return nil
}

// UpdatePupilGlobalRecord TODO: Add description
func UpdatePupilGlobalRecord(
	pupilID string, newPupil tenantmodels.Pupil, workspaceDB *sql.DB,
) error {
	var err error

	domains, err := GetAllDomainsHelper(workspaceDB)
	if err != nil {
		return err
	}

	// Check if the email ends with a valid domain
	validDomain := false
	for _, domain := range domains {
		if strings.HasSuffix(newPupil.Email, domain.Domain) {
			validDomain = true
			break
		}
	}
	if !validDomain {
		return fmt.Errorf("email učenika mora biti iz validne domene")
	}

	tx, err := workspaceDB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	accountsQuery := `UPDATE accounts a JOIN pupil_global pg ON a.id = pg.account_id
	SET a.email = ? WHERE pg.id = ? AND a.account_type = 'pupil'`

	_, err = tx.Exec(
		accountsQuery,
		newPupil.Email,
		pupilID,
	)
	if err != nil {
		if IsDuplicateEmailError(err) {
			return fmt.Errorf("korisnik sa ovim emailom već postoji")
		}
		return fmt.Errorf("error updating pupil account: %v", err)
	}

	updatePupilsQuery := `UPDATE pupil_global SET name=?, last_name=?, gender=?,
	address=?, guardian_name=?, phone_number=?, guardian_number=?, date_of_birth=?,
	religion=?, place_of_birth=? WHERE id=?`

	_, err = tx.Exec(
		updatePupilsQuery,
		newPupil.Name,
		newPupil.LastName,
		newPupil.Gender,
		newPupil.Address,
		newPupil.GuardianName,
		newPupil.PhoneNumber,
		newPupil.GuardianNumber,
		newPupil.DateOfBirth,
		newPupil.Religion,
		newPupil.PlaceOfBirth,
		pupilID,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
		}
		return fmt.Errorf("error updating pupil: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// UpdatePupilTenantRecord TODO: Add description
func UpdatePupilTenantRecord(
	pupilID string, newPupil tenantmodels.Pupil, tenantDB *sql.DB,
) error {
	updatePupilsQuery := `UPDATE pupils SET name=?, last_name=?, gender=?,
	address=?, guardian_name=?, phone_number=?, guardian_number=?, date_of_birth=?,
	religion=?, place_of_birth=? WHERE id=?`

	_, err := tenantDB.Exec(
		updatePupilsQuery,
		newPupil.Name,
		newPupil.LastName,
		newPupil.Gender,
		newPupil.Address,
		newPupil.GuardianName,
		newPupil.PhoneNumber,
		newPupil.GuardianNumber,
		newPupil.DateOfBirth,
		newPupil.Religion,
		newPupil.PlaceOfBirth,
		pupilID,
	)
	if err != nil {
		return fmt.Errorf("error updating pupil: %v", err)
	}

	return nil
}

// GetTenantsForPupil TODO: Add description
func GetTenantsForPupil(
	pupil tenantmodels.Pupil,
	db *sql.DB,
) ([]wpmodels.Tenant, error) {
	query := `SELECT tenant_id FROM pupil_tenant WHERE pupil_id = ?`
	rows, err := db.Query(query, pupil.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying pupil_tenant in util function")
	}
	defer rows.Close()

	var tenants []wpmodels.Tenant
	for rows.Next() {
		var tenantID string

		if err := rows.Scan(&tenantID); err != nil {
			return nil, fmt.Errorf("error scanning tenant_id: %v", err)
		}

		tenant, err := GetTenantByID(tenantID, db)
		if err != nil {
			return nil, fmt.Errorf(
				"error getting tenant by ID in GetTeacherTenantInstances: %v",
				err,
			)
		}

		tenants = append(tenants, *tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return tenants, nil
}

// UpdatePupilBehaviourGradeHelper updates a pupil behaviour grade for a pupil in a
// specific section
func UpdatePupilBehaviourGradeHelper(
	signature string,
	behaviourGrade tenantmodels.BehaviourGrade,
	tenantDB *sql.DB,
) (updatedBehaviourGrade *tenantmodels.BehaviourGrade, err error) {

	tx, err := tenantDB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `UPDATE pupil_behaviour SET behaviour = ?, signature = ?
	WHERE pupil_id = ? AND section_id = ? AND semester_code = ?`

	_, err = tx.Exec(
		query,
		behaviourGrade.Behaviour,
		signature,
		behaviourGrade.PupilID,
		behaviourGrade.SectionID,
		behaviourGrade.SemesterCode,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	updatedBehaviourGrade, err = GetPupilBehaviourGradeByID(
		behaviourGrade.ID, tenantDB,
	)
	if err != nil {
		return nil, err
	}

	return updatedBehaviourGrade, nil
}

// GetPupilStatisticsFieldsByPupilID retrieves the statistics fields for a pupil by their ID
func GetPupilStatisticsFieldsByPupilID(
	pupilID int,
	workspaceDB *sql.DB,
) (*tenantmodels.PupilStatistics, error) {

	var stats tenantmodels.PupilStatistics
	query := `SELECT child_of_martyr, father_name, mother_name, parents_rvi,
	living_condition, student_dorm, refugee, returnee_from_abroad,
	country_of_birth, country_of_living, citizenship,
	ethnicity, father_occupation, mother_occupation, has_no_parents,
	extra_information, child_alone, is_commuter, commuting_type,
	distance_to_school_km, has_hifz, special_honors
	FROM pupil_global WHERE id = ?`

	err := workspaceDB.QueryRow(query, pupilID).Scan(
		&stats.ChildOfMartyr, &stats.FatherName, &stats.MotherName,
		&stats.ParentsRVI, &stats.LivingCondition, &stats.StudentDorm,
		&stats.Refugee, &stats.ReturneeFromAbroad,
		&stats.CountryOfBirth, &stats.CountryOfLiving, &stats.Citizenship,
		&stats.Ethnicity, &stats.FatherOccupation, &stats.MotherOccupation,
		&stats.HasNoParents, &stats.ExtraInformation, &stats.ChildAlone,
		&stats.IsCommuter, &stats.CommutingType, &stats.DistanceToSchoolKm,
		&stats.HasHifz, &stats.SpecialHonors,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// UpdateStatisticsFieldsForPupil updates the statistics fields for a pupil
func UpdateStatisticsFieldsForPupil(
	pupilID int,
	stats *tenantmodels.PupilStatistics,
	workspaceDB *sql.DB,
) (*tenantmodels.PupilStatistics, error) {
	// Set all nil boolean fields to false
	if stats.ChildOfMartyr == nil {
		b := false
		stats.ChildOfMartyr = &b
	}
	if stats.ParentsRVI == nil {
		b := false
		stats.ParentsRVI = &b
	}
	if stats.StudentDorm == nil {
		b := false
		stats.StudentDorm = &b
	}
	if stats.Refugee == nil {
		b := false
		stats.Refugee = &b
	}
	if stats.ReturneeFromAbroad == nil {
		b := false
		stats.ReturneeFromAbroad = &b
	}
	if stats.HasNoParents == nil {
		b := false
		stats.HasNoParents = &b
	}
	if stats.IsCommuter == nil {
		b := false
		stats.IsCommuter = &b
	}
	if stats.HasHifz == nil {
		b := false
		stats.HasHifz = &b
	}
	if stats.SpecialHonors == nil {
		b := false
		stats.SpecialHonors = &b
	}

	query := `UPDATE pupil_global SET child_of_martyr = ?, father_name = ?,
	mother_name = ?, parents_rvi = ?, living_condition = ?, student_dorm = ?,
	refugee = ?, returnee_from_abroad = ?, country_of_birth = ?,
	country_of_living = ?, citizenship = ?, ethnicity = ?, father_occupation = ?,
	mother_occupation = ?, has_no_parents = ?, extra_information = ?, child_alone = ?,
	is_commuter = ?, commuting_type = ?, distance_to_school_km = ?, has_hifz = ?,
	special_honors = ? WHERE id = ?`

	_, err := workspaceDB.Exec(
		query,
		stats.ChildOfMartyr, stats.FatherName, stats.MotherName, stats.ParentsRVI,
		stats.LivingCondition, stats.StudentDorm, stats.Refugee, stats.ReturneeFromAbroad,
		stats.CountryOfBirth, stats.CountryOfLiving, stats.Citizenship,
		stats.Ethnicity, stats.FatherOccupation, stats.MotherOccupation, stats.HasNoParents,
		stats.ExtraInformation, stats.ChildAlone, stats.IsCommuter, stats.CommutingType,
		stats.DistanceToSchoolKm, stats.HasHifz, stats.SpecialHonors,
		pupilID,
	)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
