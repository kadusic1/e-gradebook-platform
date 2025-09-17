package api

import (
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/tenantfactory"
	"ednevnik-backend/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// GetTeacher returns a teacher by id (super admin only)
func GetTeacher(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	teacher, err := util.GetTeacherByID(id, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// CreateTenant creates a new tenant (super admin only)
func CreateTenant(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Decode the tenant data from the request body
	var tenant wpmodels.Tenant
	if err := json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if tenant with same email already exists
	emailDuplicateCheckQuery := `SELECT COUNT(*) FROM tenant WHERE email = ?`
	var count int
	err := userWorkspaceDb.QueryRow(emailDuplicateCheckQuery, tenant.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Institucija sa ovom email adresom već postoji.", http.StatusConflict)
		return
	}

	// Check if domain already exists
	exists, err := util.GlobalDomainExists(*tenant.Domain, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Ova domena već postoji kao globalna domena.", http.StatusConflict)
		return
	}

	// Create tenant admin
	tenantAdminData := wpmodels.Teacher{
		Name:         tenant.TeacherName,
		LastName:     tenant.TeacherLastName,
		Email:        tenant.TeacherEmail,
		Phone:        tenant.TeacherPhone,
		Password:     tenant.TeacherPassword,
		Contractions: tenant.TeacherContractions,
		Title:        tenant.TeacherTitle,
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tenantAdmin, err := util.CreateTeacher(
		tenantAdminData,
		userWorkspaceDb,
		claims.ID,
		"tenant_admin",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if *tenant.Domain == "" {
		tenant.Domain = nil
	}

	// Insert tenant data into workspace DB
	res, err := userWorkspaceDb.Exec(`INSERT INTO tenant
	(tenant_name, canton_code, address, phone, email, director_name, tenant_type,
	domain, color_config, teacher_display, teacher_invite_display, pupil_display,
	pupil_invite_display, section_display, curriculum_display, semester_display,
	lesson_display, absence_display, classroom_display, tenant_admin_id,
	tenant_city, specialization)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tenant.TenantName, tenant.CantonCode, tenant.Address, tenant.Phone,
		tenant.Email, tenant.DirectorName, tenant.TenantType, tenant.Domain,
		tenant.ColorConfig, tenant.TeacherDisplay, tenant.TeacherInviteDisplay,
		tenant.PupilDisplay, tenant.PupilInviteDisplay, tenant.SectionDisplay,
		tenant.CurriculumDisplay, tenant.SemesterDisplay, tenant.LessonDisplay,
		tenant.AbsenceDisplay, tenant.ClassroomDisplay, tenantAdmin.GetID(),
		tenant.TenantCity, tenant.Specialization,
	)
	if err != nil {
		_ = util.DeleteTeacher(
			fmt.Sprintf("%d", tenantAdmin.GetID()),
			userWorkspaceDb,
		)
		if util.IsDuplicatePhoneError(err) {
			http.Error(w, "Institucija sa ovim brojem telefona već postoji.", http.StatusConflict)
			return
		}
		if util.IsDuplicateDomain(err) {
			http.Error(w, "Institucija sa ovom domenom već postoji.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tenant.ID = id

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.CreateDB(
		tenant,
		r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create schema for the tenant database
	if err := tenantInstance.CreateSchema(*tenantAdmin); err != nil {
		http.Error(w, "DB Provision error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tenant.TeacherPassword = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// UpdateTenant updates a tenant (super admin only)
func UpdateTenant(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var err error

	// Start a transaction
	tx, err := userWorkspaceDb.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	var tenant wpmodels.Tenant
	if err = json.NewDecoder(r.Body).Decode(&tenant); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if tenant with same email already exists
	emailDuplicateCheckQuery := `SELECT COUNT(*) FROM tenant WHERE email = ?
	AND id <> ?`
	var count int
	err = tx.QueryRow(emailDuplicateCheckQuery, tenant.Email, id).Scan(&count)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Institucija sa ovom email adresom već postoji.", http.StatusConflict)
		return
	}

	// Check if domain already exists
	exists, err := util.GlobalDomainExists(*tenant.Domain, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Ova domena već postoji kao globalna domena.", http.StatusConflict)
		return
	}

	// Update tenant admin
	_, err = tx.Exec(`
    UPDATE teachers t
    JOIN accounts a ON t.account_id = a.id
    SET 
        t.name = ?,
        t.last_name = ?,
        t.phone = ?,
        a.email = ?,
		t.contractions = ?,
		t.title = ?
    WHERE t.id = (SELECT tenant_admin_id FROM tenant WHERE id = ?)`,
		tenant.TeacherName, tenant.TeacherLastName,
		tenant.TeacherPhone, tenant.TeacherEmail,
		tenant.TeacherContractions, tenant.TeacherTitle, id)
	if err != nil {
		if util.IsDuplicateEmailError(err) {
			http.Error(w, "Korisnik sa ovom email adresom već postoji.", http.StatusConflict)
			return
		}
		if util.IsDuplicatePhoneError(err) {
			http.Error(w, "Korisnik sa ovim brojem telefona već postoji.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if *tenant.Domain == "" {
		tenant.Domain = nil
	}

	_, err = tx.Exec(`UPDATE tenant SET tenant_name=?, canton_code=?,
	address=?, phone=?, email=?, director_name=?, tenant_type=?,
	domain=?, color_config=?, teacher_display=?, teacher_invite_display=?,
	pupil_display=?, pupil_invite_display=?, section_display=?,
	curriculum_display=?, semester_display=?, lesson_display=?,
	absence_display=?, classroom_display=?, tenant_city=?, specialization=?
	WHERE id=?`,
		tenant.TenantName, tenant.CantonCode, tenant.Address, tenant.Phone,
		tenant.Email, tenant.DirectorName, tenant.TenantType,
		tenant.Domain, tenant.ColorConfig, tenant.TeacherDisplay,
		tenant.TeacherInviteDisplay, tenant.PupilDisplay,
		tenant.PupilInviteDisplay, tenant.SectionDisplay,
		tenant.CurriculumDisplay, tenant.SemesterDisplay, tenant.LessonDisplay,
		tenant.AbsenceDisplay, tenant.ClassroomDisplay, tenant.TenantCity,
		tenant.Specialization, id)
	if err != nil {
		if util.IsDuplicatePhoneError(err) {
			http.Error(w, "Institucija sa ovim brojem telefona već postoji.", http.StatusConflict)
			return
		}
		if util.IsDuplicateDomain(err) {
			http.Error(w, "Institucija sa ovom domenom već postoji.", http.StatusConflict)
			return
		}
		if util.IsDuplicateEmailError(err) {
			http.Error(w, "Institucija sa ovom email adresom već postoji.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// DeleteTenant deletes a tenant (super admin only)
func DeleteTenant(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.TenantFactory(id, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Drop the tenant database
	err = tenantInstance.DropDB()
	if err != nil {
		http.Error(w, "Failed to drop tenant DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete tenant admin
	_, err = userWorkspaceDb.Exec(`
		DELETE a FROM accounts a 
		JOIN teachers t ON a.id = t.account_id 
		WHERE t.id = (SELECT tenant_admin_id FROM tenant WHERE id = ?)`, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For now the tenant deletion is cascaded with the query above
	// _, err = userWorkspaceDb.Exec(`DELETE FROM tenant WHERE id = ?`, id)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(http.StatusNoContent)
}

// ListTenants returns all tenants (super admin only)
func ListTenants(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	rows, err := userWorkspaceDb.Query(`SELECT t.id, t.tenant_name, t.canton_code,
	t.address, t.phone, t.email, t.director_name, t.tenant_type,
	t.domain, t.color_config, t.teacher_display, t.teacher_invite_display,
	t.pupil_display, t.pupil_invite_display, t.section_display,
	t.curriculum_display, t.semester_display, t.tenant_admin_id,
	tch.name, tch.last_name, tch.phone, a.email, t.lesson_display, t.absence_display,
	t.classroom_display, t.tenant_city, tch.contractions, tch.title, t.specialization
	FROM tenant t
	JOIN teachers tch ON tch.id = t.tenant_admin_id
	JOIN accounts a ON a.id = tch.account_id`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tenants []wpmodels.Tenant
	for rows.Next() {
		var s wpmodels.Tenant
		err := rows.Scan(&s.ID, &s.TenantName, &s.CantonCode, &s.Address, &s.Phone,
			&s.Email, &s.DirectorName, &s.TenantType, &s.Domain, &s.ColorConfig,
			&s.TeacherDisplay, &s.TeacherInviteDisplay, &s.PupilDisplay,
			&s.PupilInviteDisplay, &s.SectionDisplay, &s.CurriculumDisplay,
			&s.SemesterDisplay, &s.TeacherID, &s.TeacherName, &s.TeacherLastName,
			&s.TeacherPhone, &s.TeacherEmail, &s.LessonDisplay, &s.AbsenceDisplay,
			&s.ClassroomDisplay, &s.TenantCity, &s.TeacherContractions, &s.TeacherTitle,
			&s.Specialization)
		if err == nil {
			tenants = append(tenants, s)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

// GetCurriculumsForTenantAssignment returns curriculums availabe for tenant
// assignment
func GetCurriculumsForTenantAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	curriculums, err := tenantInstance.GetCurriculumsForAssignment()
	if err != nil {
		http.Error(w, "Failed to get curriculums: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(curriculums)
}

// AssignCurriculumsToTenant assigns one or more curriculums to tenant
func AssignCurriculumsToTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var curriculumCodes []string
	if err := json.NewDecoder(r.Body).Decode(&curriculumCodes); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(curriculumCodes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = tenantInstance.AssignCurriculumsToTenant(curriculumCodes)
	if err != nil {
		http.Error(w, "Failed to assign curriculums: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTenantCurriculums returns assigned tenant curriculums
func GetTenantCurriculums(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	curriculums, err := tenantInstance.GetCurriculumsForTenant()
	if err != nil {
		http.Error(w, "Failed to get curriculums: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(curriculums)
}

// UnassignCurriculumsFromTenant unassigns a curriculum from a tenant
func UnassignCurriculumsFromTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	// Get the tenant instance using the factory method
	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	curriculumCode := vars["curriculum_code"]
	if curriculumCode == "" {
		http.Error(w, "Missing curriculum code", http.StatusBadRequest)
		return
	}

	err = tenantInstance.UnassignCurriculumFromTenant(curriculumCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateNPPSemesterDates updates NPP semesters
func UpdateNPPSemesterDates(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		NPPCode      string `json:"npp_code"`
		SemesterCode string `json:"semester_code"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.NPPCode == "" || req.SemesterCode == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	updatedSemester, err := util.UpdateNPPSemesterDates(
		userWorkspaceDb, req.NPPCode, req.SemesterCode, req.StartDate, req.EndDate,
	)
	if err != nil {
		http.Error(w, "Failed to update NPP semester dates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSemester)
}

// GetAllNPPSemesters returns all NPP semesters
func GetAllNPPSemesters(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	nppSemesters, err := util.GetAllNPPSemesters(userWorkspaceDb)
	if err != nil {
		http.Error(w, "Failed to get NPP semesters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nppSemesters)
}

// GetAllDomainsHandler returns all domains (global and institution)
func GetAllDomainsHandler(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	domains, err := util.GetAllDomainsHelper(userWorkspaceDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(domains)
}

// CreateGlobalDomainHandler creates a global domain
func CreateGlobalDomainHandler(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var domain wpmodels.Domain
	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if domain.Domain == "" {
		http.Error(w, "Domain name is required", http.StatusBadRequest)
		return
	}

	err := util.InsertGlobalDomainHelper(userWorkspaceDb, domain.Domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteGlobalDomainHandler deletes a global domain
func DeleteGlobalDomainHandler(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var domain wpmodels.Domain
	if err := json.NewDecoder(r.Body).Decode(&domain); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := util.DeleteGlobalDomainHelper(userWorkspaceDb, domain.Domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
