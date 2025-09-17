package tenantmodels

import (
	"database/sql"
	"ednevnik-backend/models/interfaces"
	"fmt"
)

// Pupil represents a pupil (student) in the system
type Pupil struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	LastName         string `json:"last_name"`
	JMBG             string `json:"jmbg,omitempty"`
	Gender           string `json:"gender,omitempty"`
	Address          string `json:"address,omitempty"`
	GuardianName     string `json:"guardian_name"`
	PhoneNumber      string `json:"phone_number"`
	GuardianNumber   string `json:"guardian_number,omitempty"`
	Password         string `json:"password,omitempty"`
	DateOfBirth      string `json:"date_of_birth,omitempty"`
	Religion         string `json:"religion"`
	Email            string `json:"email"`
	PlaceOfBirth     string `json:"place_of_birth,omitempty"`
	IsCommuter       string `json:"is_commuter,omitempty"`
	Unenrolled       bool   `json:"unenrolled"`
	ParentAccessCode string `json:"parent_access_code,omitempty"`
}

// PupilStatistics represents various statistics about a pupil
type PupilStatistics struct {
	ChildOfMartyr      *bool   `json:"child_of_martyr"`
	FatherName         *string `json:"father_name"`
	MotherName         *string `json:"mother_name"`
	ParentsRVI         *bool   `json:"parents_rvi"`
	LivingCondition    *string `json:"living_condition"`
	StudentDorm        *bool   `json:"student_dorm"`
	Refugee            *bool   `json:"refugee"`
	ReturneeFromAbroad *bool   `json:"returnee_from_abroad"`
	CountryOfBirth     *string `json:"country_of_birth"`
	CountryOfLiving    *string `json:"country_of_living"`
	Citizenship        *string `json:"citizenship"`
	Ethnicity          *string `json:"ethnicity"`
	FatherOccupation   *string `json:"father_occupation"`
	MotherOccupation   *string `json:"mother_occupation"`
	HasNoParents       *bool   `json:"has_no_parents"`
	ExtraInformation   *string `json:"extra_information"`
	ChildAlone         *bool   `json:"child_alone"`
	IsCommuter         *bool   `json:"is_commuter"`
	CommutingType      *string `json:"commuting_type"`
	DistanceToSchoolKm *string `json:"distance_to_school_km"`
	HasHifz            *bool   `json:"has_hifz"`
	SpecialHonors      *bool   `json:"special_honors"`
}

// PupilSection represents the relation between pupils and sections
type PupilSection struct {
	ID        int64 `json:"id"`
	PupilID   int   `json:"pupil_id"`
	SectionID int   `json:"section_id"`
}

// StudentGrade represents a grade for a pupil
type StudentGrade struct {
	ID         int64  `json:"id"`
	PupilID    int    `json:"pupil_id"`
	SectionID  int    `json:"section_id"`
	SubjectID  int    `json:"subject_id"`
	Grade      int    `json:"grade"`
	GradeDate  string `json:"grade_date"`
	Semester   string `json:"semester"`
	TenantYear string `json:"tenant_year"`
	TeacherID  int    `json:"teacher_id"`
}

// GetID TODO: Add description
func (p Pupil) GetID() int {
	return p.ID
}

// GetName TODO: Add description
func (p Pupil) GetName() string {
	return p.Name
}

// GetLastName TODO: Add description
func (p Pupil) GetLastName() string {
	return p.LastName
}

// GetEmail TODO: Add description
func (p Pupil) GetEmail() string {
	return p.Email

}

// GetPhone TODO: Add description
func (p Pupil) GetPhone() string {
	return p.PhoneNumber
}

// GetAccountType TODO: Add description
func (p Pupil) GetAccountType(workspaceDB *sql.DB) string {
	return "pupil"
}

// GetTenantIDs TODO: Add description
func (p Pupil) GetTenantIDs(workspaceDB *sql.DB) ([]string, error) {
	var tenantIDs []string

	query := "SELECT tenant_id FROM pupil_tenant WHERE pupil_id = ?"
	rows, err := workspaceDB.Query(query, p.GetID())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tenantID string
		err := rows.Scan(&tenantID)
		if err != nil {
			return nil, err
		}
		tenantIDs = append(tenantIDs, tenantID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tenantIDs, nil
}

// GetPassword TODO: Add description
func (p Pupil) GetPassword() string {
	return p.Password
}

// GetAccountID TODO: Add description
func (p Pupil) GetAccountID(workspaceDB *sql.DB) (int, error) {
	var accountID int
	query := "SELECT account_id FROM pupil_global WHERE id = ?"
	err := workspaceDB.QueryRow(query, p.GetID()).Scan(&accountID)
	if err != nil {
		return 0, fmt.Errorf("error fetching account ID: %w", err)
	}
	return accountID, nil
}

// Ensure pupil implements user interface
var _ interfaces.User = (*Pupil)(nil)
