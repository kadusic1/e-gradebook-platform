package util

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	"ednevnik-backend/models/interfaces"
	wpmodels "ednevnik-backend/models/workspace"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// ValidateIdentifier validates MySQL identifiers
func ValidateIdentifier(identifier string) error {
	// MySQL identifier rules: 1-64 chars, alphanumeric + underscore
	// Allow dots for email addresses in usernames
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9@._-]{1,64}$`)

	if !validPattern.MatchString(identifier) {
		return fmt.Errorf("invalid identifier: contains illegal characters or exceeds length limit")
	}

	// Additional check for SQL keywords or suspicious patterns
	identifier = strings.ToUpper(identifier)
	dangerousKeywords := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
		"UNION", "OR", "AND", "WHERE", "FROM", "INTO", "VALUES", "--", "/*", "*/",
		"'", "\"", ";", "\\", "EXEC", "EXECUTE",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(identifier, keyword) {
			return fmt.Errorf("identifier contains forbidden keyword or character: %s", keyword)
		}
	}

	return nil
}

// GetTeacherByID TODO: Add description
func GetTeacherByID(id string, db interfaces.DatabaseQuerier) (wpmodels.Teacher, error) {
	query := `SELECT t.id, t.name, t.last_name, a.email, t.phone, t.contractions, t.title
	FROM
	teachers t
	JOIN accounts a ON t.account_id = a.id
	WHERE t.id = ? AND a.account_type IN ('teacher', 'tenant_admin', 'root');`
	row := db.QueryRow(query, id)

	var teacher wpmodels.Teacher
	err := row.Scan(
		&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.Email, &teacher.Phone,
		&teacher.Contractions, &teacher.Title,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return wpmodels.Teacher{}, fmt.Errorf("no teacher found with ID %s", id)
		}
		return wpmodels.Teacher{}, fmt.Errorf("error retrieving teacher: %v", err)
	}

	return teacher, nil
}

// ListTeachers TODO: Add description
func ListTeachers(
	workspaceDB *sql.DB, loggedInTeacher wpmodels.Claims,
) ([]wpmodels.Teacher, error) {
	var query string
	var rows *sql.Rows
	var err error

	if loggedInTeacher.AccountType == "root" {
		query = `SELECT t.id, t.name, t.last_name, a.email, t.phone, t.contractions, t.title
		FROM
		teachers t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.account_type='teacher';`
		rows, err = workspaceDB.Query(query)
	} else {
		query = `SELECT t.id, t.name, t.last_name, a.email, t.phone, t.contractions, t.title
		FROM
		teachers t
		JOIN accounts a ON t.account_id = a.id
		WHERE a.created_by_teacher_id = ?
		AND a.account_type='teacher';`
		rows, err = workspaceDB.Query(query, loggedInTeacher.ID)
	}

	if err != nil {
		return nil, fmt.Errorf("error listing teachers: %v", err)
	}
	defer rows.Close()

	var teachers []wpmodels.Teacher
	for rows.Next() {
		var t wpmodels.Teacher
		err := rows.Scan(
			&t.ID, &t.Name, &t.LastName, &t.Email, &t.Phone,
			&t.Contractions, &t.Title,
		)
		if err == nil {
			teachers = append(teachers, t)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return teachers, nil
}

// GetAllRegularTeachers TODO: Add description
func GetAllRegularTeachers(workspaceDB *sql.DB) ([]wpmodels.Teacher, error) {
	query := `SELECT t.id, t.name, t.last_name, a.email, t.phone FROM
	teachers t
	JOIN accounts a ON t.account_id = a.id
	WHERE a.account_type = 'teacher';`

	rows, err := workspaceDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error listing teachers: %v", err)
	}
	defer rows.Close()

	var teachers []wpmodels.Teacher
	for rows.Next() {
		var t wpmodels.Teacher
		err := rows.Scan(&t.ID, &t.Name, &t.LastName, &t.Email, &t.Phone)
		if err == nil {
			teachers = append(teachers, t)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return teachers, nil
}

// GetTeacherTablePrivileges TODO: Add description
func GetTeacherTablePrivileges() []struct {
	Table   string
	Actions string
} {
	return []struct {
		Table   string
		Actions string
	}{
		{"pupils", "SELECT, INSERT, UPDATE, DELETE"},
		{"sections", "SELECT, UPDATE"},
		{"pupils_sections", "SELECT, INSERT, UPDATE, DELETE"},
		{"student_grades", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupils_sections_invite", "SELECT, INSERT, UPDATE, DELETE"},
		{"teachers_sections_invite", "SELECT"},
		{"teachers_sections_invite_subjects", "SELECT"},
		{"homeroom_assignments", "SELECT"},
		{"teachers_sections", "SELECT"},
		{"teachers_sections_subjects", "SELECT"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_attendance", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_behaviour", "SELECT, INSERT, UPDATE, DELETE"},
	}
}

// GetServiceUserTablePrivileges TODO: Add description
func GetServiceUserTablePrivileges() []struct {
	Table   string
	Actions string
} {
	return []struct {
		Table   string
		Actions string
	}{
		{"pupils", "SELECT, INSERT, UPDATE, DELETE"},
		{"sections", "SELECT, UPDATE"},
		{"pupils_sections", "SELECT, INSERT, UPDATE, DELETE"},
		{"student_grades", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupils_sections_invite", "SELECT, INSERT, UPDATE, DELETE"},
		{"teachers_sections_invite", "SELECT, UPDATE"},
		{"teachers_sections_invite_subjects", "SELECT"},
		{"teachers_sections", "INSERT, SELECT"},
		{"teachers_sections_subjects", "INSERT, SELECT"},
		{"homeroom_assignments", "SELECT, INSERT, UPDATE"},
		{"classroom", "SELECT"},
		{"schedule", "SELECT"},
		{"time_periods", "SELECT"},
		{"class_lesson", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_attendance", "SELECT, INSERT, UPDATE, DELETE"},
		{"pupil_behaviour", "SELECT, INSERT, UPDATE, DELETE"},
	}
}

// GetTenantsForTeacher TODO: Add description
func GetTenantsForTeacher(
	teacher wpmodels.Teacher,
	db *sql.DB,
) ([]wpmodels.Tenant, error) {
	query := `SELECT tenant_id FROM teacher_tenant WHERE teacher_id = ?`
	rows, err := db.Query(query, teacher.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying teacher_tenant in GetTeacherTenantInstances util function")
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

// GetTeacherByEmail tries to find a teacher by email in the given DBs (tries each in order)
func GetTeacherByEmail(db *sql.DB, email string) (*wpmodels.Teacher, error) {
	var teacher wpmodels.Teacher
	getTeacherQuery := `SELECT t.id, t.name, t.last_name, a.email, t.phone, a.password,
	t.contractions, t.title
	FROM teachers t
	JOIN accounts a ON t.account_id = a.id
	WHERE a.email = ? AND a.account_type IN ('teacher','tenant_admin','root');`

	err := db.QueryRow(getTeacherQuery, email).Scan(
		&teacher.ID,
		&teacher.Name,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Phone,
		&teacher.Password,
		&teacher.Contractions,
		&teacher.Title,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &teacher, nil
}

// CreateTeacher TODO: Add description
func CreateTeacher(
	teacher wpmodels.Teacher,
	workspaceDB *sql.DB,
	loggedInUserID int,
	accountType string,
) (*wpmodels.Teacher, error) {
	var err error
	if err = ValidateIdentifier(teacher.Email); err != nil {
		return nil, fmt.Errorf("ovaj email nije validan")
	}

	tx, err := workspaceDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	hash, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}
	teacher.Password = string(hash)

	// First insert account data
	accountInsertQuery := `INSERT INTO accounts (email, password,
	created_by_teacher_id, account_type) VALUES (?, ?, ?, ?)`

	res, err := tx.Exec(
		accountInsertQuery,
		teacher.Email,
		teacher.Password,
		loggedInUserID,
		accountType,
	)
	if err != nil {
		if IsDuplicateEmailError(err) {
			return nil, fmt.Errorf("korisnik sa ovim emailom već postoji")
		}
		return nil, fmt.Errorf("error inserting account data: %v", err)
	}

	// Get the last inserted account ID
	accountID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Insert teacher data
	insertTeacherQuery := `INSERT INTO teachers (name, last_name, phone, account_id,
	contractions, title)
	VALUES (?, ?, ?, ?, ?, ?)`
	res, err = tx.Exec(
		insertTeacherQuery,
		teacher.Name,
		teacher.LastName,
		teacher.Phone,
		accountID,
		teacher.Contractions,
		teacher.Title,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return nil, fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
		}
		return nil, fmt.Errorf("error inserting teacher: %v", err)
	}
	teacherID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %v", err)
	}
	teacher.ID = int(teacherID)

	if err != nil {
		return nil, fmt.Errorf("error inserting teacher data: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	// Never return the password in the response
	teacher.Password = "" // Clear the password before returning

	return &teacher, nil
}

// RegisterTeacher TODO: Add description
func RegisterTeacher(
	teacher wpmodels.Teacher,
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
		if strings.HasSuffix(teacher.Email, domain.Domain) {
			validDomain = true
			break
		}
	}
	if !validDomain {
		return fmt.Errorf("email nastavnika mora biti iz validne domene")
	}

	exists, err := AccountWithEmailExists(teacher.Email, workspaceDB)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("korisnik sa ovim emailom već postoji")
	}

	exists, err = TeacherWithPhoneExists(teacher.Phone, workspaceDB)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
	}

	if err = ValidateIdentifier(teacher.Email); err != nil {
		return fmt.Errorf("ovaj email nije validan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(teacher.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	teacher.Password = string(hash)

	tx, err := workspaceDB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// First insert account data
	accountInsertQuery := `INSERT INTO pending_accounts (email, password, account_type)
	VALUES (?, ?, 'teacher')`

	res, err := tx.Exec(
		accountInsertQuery,
		teacher.Email,
		teacher.Password,
	)
	if err != nil {
		if IsDuplicateEmailError(err) {
			return fmt.Errorf(
				"neverifikovani korisnik sa ovim emailom već postoji - provjerite email za aktivaciju",
			)
		}
		return fmt.Errorf("error inserting account data: %v", err)
	}

	// Get the last inserted account ID
	accountID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	// Insert teacher data
	insertTeacherQuery := `INSERT INTO pending_teachers (name, last_name, phone, account_id,
		contractions, title)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(
		insertTeacherQuery,
		teacher.Name,
		teacher.LastName,
		teacher.Phone,
		accountID,
		teacher.Contractions,
		teacher.Title,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return fmt.Errorf(
				"neverifikovani korisnik sa ovim brojem telefona već postoji - provjerite email za aktivaciju",
			)
		}
		return fmt.Errorf("error inserting teacher data: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	verificationToken, err := GetPendingAccountVerificationToken(int(accountID), workspaceDB)
	if err != nil {
		return fmt.Errorf("error getting verification token: %v", err)
	}

	go func() {
		_ = SendVerificationEmail(
			teacher.Email,
			fmt.Sprintf("%s %s", teacher.Name, teacher.LastName),
			fmt.Sprintf("%s/verify?token=%s", os.Getenv("FRONTEND_URL"), verificationToken),
		)
	}()

	return nil
}

// UpdateTeacher TODO: Add description
func UpdateTeacher(
	teacher wpmodels.Teacher,
	oldTeacher wpmodels.Teacher,
	workspaceDB *sql.DB,
) (*wpmodels.Teacher, error) {
	var err error

	domains, err := GetAllDomainsHelper(workspaceDB)
	if err != nil {
		return nil, err
	}

	// Check if the email ends with a valid domain
	validDomain := false
	for _, domain := range domains {
		if strings.HasSuffix(teacher.Email, domain.Domain) {
			validDomain = true
			break
		}
	}
	if !validDomain {
		return nil, fmt.Errorf("email nastavnika mora biti iz validne domene")
	}

	if err = ValidateIdentifier(teacher.Email); err != nil {
		return nil, fmt.Errorf("ovaj email nije validan")
	}

	// Start a transaction to ensure atomicity
	tx, err := workspaceDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if teacher.Email != oldTeacher.Email {

		accountsQuery := `UPDATE accounts a
		JOIN teachers t ON a.id = t.account_id
		SET a.email = ? WHERE t.id = ?`
		_, err = tx.Exec(
			accountsQuery,
			teacher.Email,
			oldTeacher.ID,
		)
		if err != nil {
			if IsDuplicateEmailError(err) {
				return nil, fmt.Errorf("korisnik sa ovim emailom već postoji")
			}
			return nil, fmt.Errorf("error updating teacher email: %v", err)
		}
	}

	updateQuery := `UPDATE teachers SET name = ?, last_name = ?, phone = ?,
	contractions = ?, title = ? WHERE id = ?`
	_, err = tx.Exec(
		updateQuery,
		teacher.Name,
		teacher.LastName,
		teacher.Phone,
		teacher.Contractions,
		teacher.Title,
		teacher.ID,
	)
	if err != nil {
		if IsDuplicatePhoneError(err) {
			return nil, fmt.Errorf("korisnik sa ovim brojem telefona već postoji")
		}
		return nil, fmt.Errorf("error updating teacher: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error committing transaction: %v", err)
	}

	return &teacher, nil
}

// DeleteTeacher TODO: Add description
func DeleteTeacher(
	teacherID string,
	workspaceDB *sql.DB,
) error {
	teacher, err := GetTeacherByID(teacherID, workspaceDB)
	if err != nil {
		return fmt.Errorf("error getting teacher by ID: %v", err)
	}

	teacherAccountID, err := teacher.GetAccountID(workspaceDB)
	if err != nil {
		return fmt.Errorf("error getting teacher account ID: %v", err)
	}

	deleteQuery := `DELETE FROM accounts WHERE id = ?;`

	_, err = workspaceDB.Exec(deleteQuery, teacherAccountID)
	if err != nil {
		return fmt.Errorf("error deleting teacher: %v", err)
	}
	return nil
}

// GetAllAssignedSubjectsMap TODO: Add description
func GetAllAssignedSubjectsMap(tenantDB *sql.DB) (map[int]map[int][]wpmodels.Subject, error) {
	query := `
        SELECT tss.teacher_id, tss.section_id, s.subject_code, s.subject_name
        FROM ednevnik_workspace.subjects s
        JOIN teachers_sections_subjects tss ON s.subject_code = tss.subject_code`
	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[int][]wpmodels.Subject)
	for rows.Next() {
		var teacherID, sectionID int
		var s wpmodels.Subject
		if err := rows.Scan(&teacherID, &sectionID, &s.SubjectCode, &s.SubjectName); err != nil {
			return nil, err
		}
		if _, ok := result[teacherID]; !ok {
			result[teacherID] = make(map[int][]wpmodels.Subject)
		}
		result[teacherID][sectionID] = append(result[teacherID][sectionID], s)
	}
	return result, nil
}

// GetAllPendingSubjectsMap performs a Bulk fetch: All pending subjects for all
// teacher-section pairs Map where key is teacher ID and value is another map
// where key is section ID and value is a slice of subjects
func GetAllPendingSubjectsMap(tenantDB *sql.DB) (map[int]map[int][]wpmodels.Subject, error) {
	query := `
        SELECT tsi.teacher_id, tsi.section_id, s.subject_code, s.subject_name
        FROM ednevnik_workspace.subjects s
        JOIN teachers_sections_invite_subjects tsis ON s.subject_code = tsis.subject_code
        JOIN teachers_sections_invite tsi ON tsi.id = tsis.invite_id
        WHERE tsi.status = 'pending'`
	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[int][]wpmodels.Subject)
	for rows.Next() {
		var teacherID, sectionID int
		var s wpmodels.Subject
		if err := rows.Scan(&teacherID, &sectionID, &s.SubjectCode, &s.SubjectName); err != nil {
			return nil, err
		}
		if _, ok := result[teacherID]; !ok {
			result[teacherID] = make(map[int][]wpmodels.Subject)
		}
		result[teacherID][sectionID] = append(result[teacherID][sectionID], s)
	}
	return result, nil
}

// GetAllPendingInviteIDsMap returns a map where key is teacher ID and value is
// another map where key is section ID and value pending inviteID
func GetAllPendingInviteIDsMap(tenantDB *sql.DB) (map[int]map[int]int, error) {
	query := `
        SELECT tsi.teacher_id, tsi.section_id, tsi.id as invite_id
        FROM teachers_sections_invite tsi
        WHERE tsi.status = 'pending'`

	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[int]int)
	for rows.Next() {
		var teacherID, sectionID, inviteID int
		if err := rows.Scan(&teacherID, &sectionID, &inviteID); err != nil {
			return nil, err
		}

		// Check if teacher exists in the map, if not create the inner map
		if _, exists := result[teacherID]; !exists {
			result[teacherID] = make(map[int]int)
		}

		// Add the section ID and invite ID to the teacher's map
		result[teacherID][sectionID] = inviteID
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAvailableSubjects returns subjects from allSubjects that are not in assigned or pending
func GetAvailableSubjects(
	allSubjects, assignedSubjects, pendingSubjects []wpmodels.Subject,
) []wpmodels.Subject {
	assigned := make(map[string]struct{})
	for _, s := range assignedSubjects {
		assigned[s.SubjectCode] = struct{}{}
	}
	for _, s := range pendingSubjects {
		assigned[s.SubjectCode] = struct{}{}
	}
	var available []wpmodels.Subject
	for _, s := range allSubjects {
		if _, found := assigned[s.SubjectCode]; !found {
			available = append(available, s)
		}
	}
	return available
}

// MakeHomeroomTeachersMap TODO: Add description
func MakeHomeroomTeachersMap(rows *sql.Rows) (map[int]int, error) {
	result := make(map[int]int)

	for rows.Next() {
		var sectionID int
		var teacherID int

		if err := rows.Scan(&sectionID, &teacherID); err != nil {
			return nil, fmt.Errorf("error scanning rows: %v", err)
		}

		result[sectionID] = teacherID
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return result, nil
}

// GetAllHomeroomTeachersMapHelper returns a map of section ID key and teacher ID value
func GetAllHomeroomTeachersMapHelper(tenantDB *sql.DB) (
	map[int]int, error,
) {
	query := `SELECT section_id, teacher_id FROM homeroom_assignments`

	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying homeroom assignments")
	}
	defer rows.Close()

	result, err := MakeHomeroomTeachersMap(rows)
	if err != nil {
		return nil, fmt.Errorf("error making section teacher map: %v", err)
	}

	return result, nil
}

// GetAllPendingHomeroomTeachersMapHelper returns a map of section ID key and
// teacher ID value
func GetAllPendingHomeroomTeachersMapHelper(tenantDB *sql.DB) (
	map[int]int, error,
) {
	query := `SELECT section_id, teacher_id FROM teachers_sections_invite
	WHERE homeroom_teacher = 1 AND status = 'pending'`

	rows, err := tenantDB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying homeroom assignments")
	}
	defer rows.Close()

	result, err := MakeHomeroomTeachersMap(rows)
	if err != nil {
		return nil, fmt.Errorf("error making section teacher map: %v", err)
	}

	return result, nil
}

// GetDataForTeacherInviteForTenantHelper TODO: Add description
func GetDataForTeacherInviteForTenantHelper(
	tenantID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) ([]commonmodels.DataForTeacherSectionInvite, error) {

	sections, err := ListSectionsForTenant(
		fmt.Sprintf("%d", tenantID), 0, tenantDB, workspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing sections for tenant: %v", err)
	}

	teachers, err := GetAllRegularTeachers(workspaceDB)
	if err != nil {
		return nil, fmt.Errorf("error getting all regular teachers: %v", err)
	}

	assignedMap, err := GetAllAssignedSubjectsMap(tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error fetching assigned subjects: %v", err)
	}
	pendingMap, err := GetAllPendingSubjectsMap(tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error fetching pending subjects: %v", err)
	}

	pendingInvitesMap, err := GetAllPendingInviteIDsMap(tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error fetching pending invite IDs: %v", err)
	}

	homeroomAssignments, err := GetAllHomeroomTeachersMapHelper(tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error getting homeroom assignments")
	}

	pendingHomeroomAssignments, err := GetAllPendingHomeroomTeachersMapHelper(tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error getting homeroom assignments")
	}

	var DataForInvite []commonmodels.DataForTeacherSectionInvite

	for _, section := range sections {
		allSubjects, err := GetAllSubjectsForCurriculumCode(
			section.CurriculumCode, workspaceDB,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting subjects for curriculum code: %v", err)
		}

		for _, teacher := range teachers {
			var inviteData commonmodels.DataForTeacherSectionInvite
			inviteData.Teacher = teacher
			inviteData.Section = section
			inviteData.AllSubjects = allSubjects

			assignedSubjects := []wpmodels.Subject{}
			if sectionSubjectsMap, teacherHasAssignments := assignedMap[teacher.ID]; teacherHasAssignments {
				if subjectsForSection, sectionHasAssignments := sectionSubjectsMap[int(section.ID)]; sectionHasAssignments {
					assignedSubjects = subjectsForSection
				}
			}
			inviteData.AssignedSubjects = assignedSubjects

			pendingSubjects := []wpmodels.Subject{}
			if sectionPendingMap, teacherHasPending := pendingMap[teacher.ID]; teacherHasPending {
				if subjectsPendingForSection, sectionHasPending := sectionPendingMap[int(section.ID)]; sectionHasPending {
					pendingSubjects = subjectsPendingForSection
				}
			}
			inviteData.PendingSubjects = pendingSubjects

			availableSubjects := GetAvailableSubjects(
				allSubjects, assignedSubjects, pendingSubjects,
			)
			inviteData.AvailableSubjects = availableSubjects

			inviteData.IsHomeroomTeacher = homeroomAssignments[int(section.ID)] == teacher.ID
			inviteData.PendingHomeroomTeacher = pendingHomeroomAssignments[int(section.ID)] == teacher.ID
			inviteData.InviteIndexID = pendingInvitesMap[teacher.ID][int(section.ID)]

			DataForInvite = append(DataForInvite, inviteData)
		}
	}

	return DataForInvite, nil
}

// GetAssignedSubjectsForTeacher fetches assigned subjects for a specific teacher
func GetAssignedSubjectsForTeacher(teacherID string, tenantDB *sql.DB) (map[int][]wpmodels.Subject, error) {
	query := `
        SELECT tss.section_id, s.subject_code, s.subject_name
        FROM ednevnik_workspace.subjects s
        JOIN teachers_sections_subjects tss ON s.subject_code = tss.subject_code
        WHERE tss.teacher_id = ?`
	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]wpmodels.Subject)
	for rows.Next() {
		var sectionID int
		var s wpmodels.Subject
		if err := rows.Scan(&sectionID, &s.SubjectCode, &s.SubjectName); err != nil {
			return nil, err
		}
		result[sectionID] = append(result[sectionID], s)
	}
	return result, nil
}

// GetPendingSubjectsForTeacher fetches pending subjects for a specific teacher
func GetPendingSubjectsForTeacher(teacherID string, tenantDB *sql.DB) (map[int][]wpmodels.Subject, error) {
	query := `
        SELECT tsi.section_id, s.subject_code, s.subject_name, tsi.id as invite_id
        FROM ednevnik_workspace.subjects s
        JOIN teachers_sections_invite_subjects tsis ON s.subject_code = tsis.subject_code
        JOIN teachers_sections_invite tsi ON tsi.id = tsis.invite_id
        WHERE tsi.teacher_id = ? AND tsi.status = 'pending'`
	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int][]wpmodels.Subject)
	for rows.Next() {
		var sectionID int
		var s wpmodels.Subject
		if err := rows.Scan(&sectionID, &s.SubjectCode, &s.SubjectName, &s.InviteIndexID); err != nil {
			return nil, err
		}
		result[sectionID] = append(result[sectionID], s)
	}
	return result, nil
}

// GetPendingInviteIDsMapForTeacher returns a map where key is teacher ID and
// value is another map where key is section ID and value pending inviteID
func GetPendingInviteIDsMapForTeacher(tenantDB *sql.DB, teacherID string) (map[int]map[int]int, error) {
	query := `
        SELECT tsi.teacher_id, tsi.section_id, tsi.id as invite_id
        FROM teachers_sections_invite tsi
        WHERE tsi.status = 'pending' AND tsi.teacher_id = ?`

	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]map[int]int)
	for rows.Next() {
		var teacherID, sectionID, inviteID int
		if err := rows.Scan(&teacherID, &sectionID, &inviteID); err != nil {
			return nil, err
		}

		// Check if teacher exists in the map, if not create the inner map
		if _, exists := result[teacherID]; !exists {
			result[teacherID] = make(map[int]int)
		}

		// Add the section ID and invite ID to the teacher's map
		result[teacherID][sectionID] = inviteID
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetHomeroomTeachersMapForTeacherHelper returns a map of section ID key and
// teacher ID value
func GetHomeroomTeachersMapForTeacherHelper(tenantDB *sql.DB, teacherID string) (
	map[int]int, error,
) {
	query := `SELECT section_id, teacher_id FROM homeroom_assignments
	WHERE teacher_id = ?`

	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error querying homeroom assignments")
	}
	defer rows.Close()

	result, err := MakeHomeroomTeachersMap(rows)
	if err != nil {
		return nil, fmt.Errorf("error making section teacher map: %v", err)
	}

	return result, nil
}

// GetPendingHomeroomTeachersMapForTeacherHelper returns a map of section ID key
// and teacher ID value
func GetPendingHomeroomTeachersMapForTeacherHelper(tenantDB *sql.DB, teacherID string) (
	map[int]int, error,
) {
	query := `SELECT section_id, teacher_id FROM teachers_sections_invite
	WHERE homeroom_teacher = 1 AND status = 'pending' AND teacher_id = ?`

	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error querying homeroom assignments")
	}
	defer rows.Close()

	result, err := MakeHomeroomTeachersMap(rows)
	if err != nil {
		return nil, fmt.Errorf("error making section teacher map: %v", err)
	}

	return result, nil
}

// GetDataForTeacherInviteForSingleTeacher returns invite data for a specific teacher
func GetDataForTeacherInviteForSingleTeacher(
	teacherID string,
	tenantID int,
	tenantDB *sql.DB,
	workspaceDB *sql.DB,
) ([]commonmodels.DataForTeacherSectionInvite, error) {

	// Get the specific teacher
	teacher, err := GetTeacherByID(teacherID, workspaceDB)
	if err != nil {
		return nil, fmt.Errorf("error getting teacher: %v", err)
	}

	// Get all sections for the tenant
	sections, err := ListSectionsForTenant(
		fmt.Sprintf("%d", tenantID), 0, tenantDB, workspaceDB,
	)
	if err != nil {
		return nil, fmt.Errorf("error listing sections for tenant: %v", err)
	}

	// Get assigned and pending subjects for this teacher
	assignedMap, err := GetAssignedSubjectsForTeacher(teacherID, tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error fetching assigned subjects: %v", err)
	}
	pendingMap, err := GetPendingSubjectsForTeacher(teacherID, tenantDB)
	if err != nil {
		return nil, fmt.Errorf("error fetching pending subjects: %v", err)
	}

	pendingInvitesMap, err := GetPendingInviteIDsMapForTeacher(tenantDB, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error fetching pending invite IDs: %v", err)
	}

	homeroomAssignments, err := GetHomeroomTeachersMapForTeacherHelper(tenantDB, teacherID)
	if err != nil {
		return nil, fmt.Errorf("error getting homeroom assignments")
	}

	pendingHomeroomAssignments, err := GetPendingHomeroomTeachersMapForTeacherHelper(
		tenantDB, teacherID,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting homeroom assignments")
	}

	var DataForInvite []commonmodels.DataForTeacherSectionInvite

	for _, section := range sections {
		allSubjects, err := GetAllSubjectsForCurriculumCode(
			section.CurriculumCode, workspaceDB,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting subjects for curriculum code: %v", err)
		}

		var inviteData commonmodels.DataForTeacherSectionInvite
		inviteData.Teacher = teacher
		inviteData.Section = section
		inviteData.AllSubjects = allSubjects

		// Get assigned subjects for this section
		assignedSubjects := []wpmodels.Subject{}
		if subjectsForSection, exists := assignedMap[int(section.ID)]; exists {
			assignedSubjects = subjectsForSection
		}
		inviteData.AssignedSubjects = assignedSubjects

		// Get pending subjects for this section
		pendingSubjects := []wpmodels.Subject{}
		if subjectsPendingForSection, exists := pendingMap[int(section.ID)]; exists {
			pendingSubjects = subjectsPendingForSection
		}
		inviteData.PendingSubjects = pendingSubjects

		// Calculate available subjects
		availableSubjects := GetAvailableSubjects(
			allSubjects, assignedSubjects, pendingSubjects,
		)
		inviteData.AvailableSubjects = availableSubjects

		inviteData.IsHomeroomTeacher = homeroomAssignments[int(section.ID)] == teacher.ID
		inviteData.PendingHomeroomTeacher = pendingHomeroomAssignments[int(section.ID)] == teacher.ID
		inviteData.InviteIndexID = pendingInvitesMap[teacher.ID][int(section.ID)]

		DataForInvite = append(DataForInvite, inviteData)
	}

	return DataForInvite, nil
}

// GetTenantIDForTenantAdmin TODO: Add description
func GetTenantIDForTenantAdmin(
	tenantAdminID int, workspaceDB *sql.DB,
) (int, error) {
	query := `SELECT id FROM tenant WHERE tenant_admin_id = ?`
	var tenantID int
	err := workspaceDB.QueryRow(query, tenantAdminID).Scan(&tenantID)
	if err != nil {
		return 0, fmt.Errorf("error retrieving tenant ID: %v", err)
	}
	return tenantID, nil
}

// GetSectionSubjectsForTeacher fetches section subjects assigned to a specific
// teacher
func GetSectionSubjectsForTeacher(
	teacherID, sectionID int, tenantDB *sql.DB,
) ([]wpmodels.Subject, error) {
	query := `
        SELECT s.subject_code, s.subject_name
        FROM ednevnik_workspace.subjects s
        JOIN teachers_sections_subjects tss ON s.subject_code = tss.subject_code
        WHERE tss.teacher_id = ? AND tss.section_id = ?`
	rows, err := tenantDB.Query(query, teacherID, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []wpmodels.Subject
	for rows.Next() {
		var s wpmodels.Subject
		if err := rows.Scan(&s.SubjectCode, &s.SubjectName); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}
