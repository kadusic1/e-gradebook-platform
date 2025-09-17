package api

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/tenantfactory"
	"ednevnik-backend/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetMetadataForSectionCreation TODO: Add description
func GetMetadataForSectionCreation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	sectionMetadata, err := tenantInstance.GetMetadataForSectionCreation()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sectionMetadata)
}

// CreateSection TODO: Add description
func CreateSection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	var section tenantmodels.SectionCreate
	if err := json.NewDecoder(r.Body).Decode(&section); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdSection, err := tenantInstance.CreateTenantSection(section)
	if err != nil {
		if util.InvalidSectionYearHelper(err) {
			http.Error(
				w,
				"Školska godina mora biti u formatu YYYY/YYYY (npr. 2025/2026)",
				http.StatusBadRequest,
			)
			return
		}
		if util.DuplicateSectionHelper(err) {
			http.Error(
				w,
				"Odjeljenje već postoji za ovaj razred i ovu školsku godinu",
				http.StatusBadRequest,
			)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdSection)
}

// GetSectionsForTenant TODO: Add description
func GetSectionsForTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	archived := vars["archived"]
	if archived == "" {
		http.Error(w, "Missing archived parameter", http.StatusBadRequest)
		return
	}

	archivedInt, err := strconv.Atoi(archived)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sections, err := tenantInstance.GetSectionsForTenant(archivedInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sections)
}

// DeleteSection TODO: Add description
func DeleteSection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Missing section_id", http.StatusBadRequest)
		return
	}
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}
	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.DeleteTenantSection(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSection TODO: Add description
func UpdateSection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Missing section_id", http.StatusBadRequest)
		return
	}

	var sectionUpdate tenantmodels.Section
	if err := json.NewDecoder(r.Body).Decode(&sectionUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(
		fmt.Sprintf("%d", sectionUpdate.TenantID),
		r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSection, err := tenantInstance.UpdateTenantSection(
		sectionUpdate, sectionID,
	)
	if err != nil {
		if util.InvalidSectionYearHelper(err) {
			http.Error(
				w,
				"Školska godina mora biti u formatu YYYY/YYYY (npr. 2025/2026)",
				http.StatusBadRequest,
			)
			return
		}
		if util.DuplicateSectionHelper(err) {
			http.Error(
				w,
				"Odjeljenje već postoji za ovaj razred i ovu školsku godinu",
				http.StatusBadRequest,
			)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newSection)
}

// ListPupilAccounts TODO: Add description
func ListPupilAccounts(w http.ResponseWriter, r *http.Request) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupils, err := util.ListPupilAccounts(workspaceDB, *claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pupils)
}

// CreatePupil TODO: Add description
func CreatePupil(w http.ResponseWriter, r *http.Request) {
	var pupil tenantmodels.Pupil
	if err := json.NewDecoder(r.Body).Decode(&pupil); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createdPupil, err := util.CreatePupil(
		pupil, workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPupil)
}

// DeletePupil TODO: Add description
func DeletePupil(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Missing pupil_id", http.StatusBadRequest)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupil, err := util.GetGlobalPupilByID(
		pupilID, workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pupilTenantIDs, err := pupil.GetTenantIDs(workspaceDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, tenantID := range pupilTenantIDs {
		tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tenantInstance.DeletePupilFromTenant(pupilID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = util.DeleteGlobalPupilRecord(
		*pupil, workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdatePupil TODO: Add description
func UpdatePupil(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Missing pupil_id", http.StatusBadRequest)
		return
	}

	var updatedPupil tenantmodels.Pupil
	if err := json.NewDecoder(r.Body).Decode(&updatedPupil); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	oldPupil, err := util.GetGlobalPupilByID(
		pupilID, workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// First update global settings
	err = util.UpdatePupilGlobalRecord(
		fmt.Sprintf("%d", oldPupil.ID),
		updatedPupil,
		workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update tenant settings
	pupilTenantIDs, err := updatedPupil.GetTenantIDs(workspaceDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, tenantID := range pupilTenantIDs {
		tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tenantInstance.UpdatePupil(*oldPupil, updatedPupil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedPupil)
}

// UpdateTenantSemesterDates TODO: Add description
func UpdateTenantSemesterDates(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TenantID     int    `json:"tenant_id"`
		SemesterCode string `json:"semester_code"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		NPPCode      string `json:"npp_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.SemesterCode == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(fmt.Sprintf("%d", req.TenantID), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedSemester, err := tenantInstance.UpdateTenantSemesterDates(
		req.SemesterCode, req.StartDate, req.EndDate, req.NPPCode,
	)
	if err != nil {
		http.Error(w, "Failed to update semester: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSemester)
}

// GetTenantSemesters TODO: Add description
func GetTenantSemesters(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenantSemesters, err := tenantInstance.GetSemestersForTenant()
	if err != nil {
		http.Error(w, "Failed to get semesters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenantSemesters)
}

// SendPupilSectionInvite TODO: Add description
func SendPupilSectionInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Missing section_id", http.StatusBadRequest)
		return
	}

	var pupilIDs []int
	if err := json.NewDecoder(r.Body).Decode(&pupilIDs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(
		tenantID, r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var invites []tenantmodels.PupilSectionInvite

	for _, pupilID := range pupilIDs {
		newInvite, err := tenantInstance.SendPupilSectionInvite(
			pupilID, sectionID, tenantID,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send invite for pupil %d: %s", pupilID, err.Error()), http.StatusInternalServerError)
			return
		}
		invites = append(invites, *newInvite)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invites)
}

// GetTeacherInviteDataForTenantHandler TODO: Add description
func GetTeacherInviteDataForTenantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := tenantInstance.GetDataForTeacherInviteForTenant()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleTeacherSectionAssignments TODO: Add description
func HandleTeacherSectionAssignments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	teacherIDStr := vars["teacher_id"]
	if teacherIDStr == "" {
		http.Error(w, "Missing teacher_id", http.StatusBadRequest)
		return
	}

	teacherID, err := strconv.Atoi(teacherIDStr)
	if err != nil {
		http.Error(w, "Invalid teacher_id", http.StatusBadRequest)
		return
	}

	var assignmentRequest wpmodels.TeacherSectionAssignment
	if err := json.NewDecoder(r.Body).Decode(&assignmentRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newData, newTeachers, err := tenantInstance.HandleTeacherSectionAssignments(
		teacherID, assignmentRequest,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type dataResponse struct {
		Data     []commonmodels.DataForTeacherSectionInvite `json:"invite_data"`
		Teachers []wpmodels.Teacher                         `json:"teacher_data"`
	}

	response := dataResponse{
		Data:     newData,
		Teachers: newTeachers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTeacherInvitesForTenantHandler TODO: Add description
func GetTeacherInvitesForTenantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	invites, err := tenantInstance.GetAllTeacherInvites()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invites)
}

// CreateTeacher creates a new teacher (super admin only)
func CreateTeacher(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var teacher wpmodels.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if teacher.Password == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	createdTeacher, err := util.CreateTeacher(
		teacher,
		userWorkspaceDb,
		claims.ID,
		"teacher",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTeacher)
}

// ListTeachers returns all teachers (super admin only)
func ListTeachers(w http.ResponseWriter, r *http.Request) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teachers, err := util.ListTeachers(
		userWorkspaceDb,
		*claims,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// UpdateTeacher updates a teacher (super admin only)
func UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	// Prevent update for teacher if he is not updating his account
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.AccountType == "teacher" && fmt.Sprintf("%d", claims.ID) != id {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userWorkspaceDb *sql.DB
	var err error
	if claims.AccountType == "teacher" {
		userWorkspaceDb, err = util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		userWorkspaceDb, ok = util.GetUserWorkspaceDBFromContext(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var oldTeacher wpmodels.Teacher
	oldTeacher, err = util.GetTeacherByID(id, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	var teacher wpmodels.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := util.UpdateTeacher(
		teacher,
		oldTeacher,
		userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

// DeleteTeacher deletes a teacher (super admin only)
func DeleteTeacher(w http.ResponseWriter, r *http.Request) {
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

	teacherToDelete, err := util.GetTeacherByID(id, userWorkspaceDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenants, err := util.GetTenantsForTeacher(teacherToDelete, userWorkspaceDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, tenant := range tenants {
		tenantInstance, err := tenantfactory.Struct(
			tenant, r,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tenantInstance.DeleteTenantTeacherData(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = util.DeleteTeacher(id, userWorkspaceDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTeachersPerTenant TODO: Add description
func ListTeachersPerTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	var teachers []wpmodels.Teacher

	teachers, err = tenantInstance.GetTeachersForTenant()
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// DeleteTeacherInviteHandler TODO: Add description
func DeleteTeacherInviteHandler(w http.ResponseWriter, r *http.Request) {
	type DeleteRequest struct {
		InviteID  int `json:"invite_id"`
		TeacherID int `json:"teacher_id"`
		TenantID  int `json:"tenant_id"`
	}

	var request DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(
		fmt.Sprintf("%d", request.TenantID), r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.DeleteTeacherInvite(
		request.InviteID,
		request.TeacherID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedInviteData, err := tenantInstance.GetDataForTeacherInviteForTenant()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInviteData)
}

// DeleteTeacherForTenantHandler TODO: Add description
func DeleteTeacherForTenantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["teacher_id"]
	if teacherID == "" {
		http.Error(w, "Missing teacher id", http.StatusBadRequest)
		return
	}

	tenantID := vars["tenant_id"]
	if teacherID == "" {
		http.Error(w, "Missing tenant id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.DeleteTeacherFromTenant(teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTenant returns a tenant by id (super admin only)
func GetTenant(w http.ResponseWriter, r *http.Request) {
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
	var tenant wpmodels.Tenant
	err := userWorkspaceDb.QueryRow(`SELECT t.id, t.tenant_name, t.canton_code, t.address,
	t.phone, t.email, t.director_name, t.tenant_type, t.domain, t.color_config,
	t.teacher_display, t.teacher_invite_display, t.pupil_display, t.pupil_invite_display,
	t.section_display, t.curriculum_display, t.semester_display, t.tenant_admin_id,
	tch.name, tch.last_name, tch.phone, a.email, t.lesson_display, t.absence_display,
	t.classroom_display, t.tenant_city
	FROM tenant t
	JOIN teachers tch ON tch.id = t.tenant_admin_id
	JOIN accounts a ON a.id = tch.account_id
	WHERE t.id = ?`, id).Scan(
		&tenant.ID, &tenant.TenantName, &tenant.CantonCode, &tenant.Address,
		&tenant.Phone, &tenant.Email, &tenant.DirectorName, &tenant.TenantType,
		&tenant.Domain, &tenant.ColorConfig, &tenant.TeacherDisplay,
		&tenant.TeacherInviteDisplay, &tenant.PupilDisplay, &tenant.PupilInviteDisplay,
		&tenant.SectionDisplay, &tenant.CurriculumDisplay, &tenant.SemesterDisplay,
		&tenant.TeacherID, &tenant.TeacherName, &tenant.TeacherLastName,
		&tenant.TeacherPhone, &tenant.TeacherEmail, &tenant.LessonDisplay,
		&tenant.AbsenceDisplay, &tenant.ClassroomDisplay, &tenant.TenantCity,
	)
	if err != nil {
		http.Error(w, "Tenant not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenant)
}

// CreateScheduleHandler TODO: Add description
func CreateScheduleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Missing section_id", http.StatusBadRequest)
		return
	}

	var data tenantmodels.ScheduleGroupCollection
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.CreateSchedule(data, sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetScheduleForSectionHandler TODO: Add description
func GetScheduleForSectionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Missing section_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schedule, err := tenantInstance.GetScheduleForSection(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

// CreateClassroomHandler TODO: Add description
func CreateClassroomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	var classroom tenantmodels.Classroom
	if err := json.NewDecoder(r.Body).Decode(&classroom); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.CreateClassroom(classroom)
	if err != nil {
		if util.DuplicatePrimaryKeyHelper(err) {
			http.Error(w, "Učionica sa ovim brojem učionice već postoji.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateClassroomHandler TODO: Add description
func UpdateClassroomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	oldCode := vars["classroom_code"]
	if oldCode == "" {
		http.Error(w, "Missing classroom_code", http.StatusBadRequest)
		return
	}

	var classroom tenantmodels.Classroom
	if err := json.NewDecoder(r.Body).Decode(&classroom); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.UpdateClassroom(classroom, oldCode)
	if err != nil {
		if util.DuplicatePrimaryKeyHelper(err) {
			http.Error(w, "Učionica sa ovim brojem učionice već postoji.", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteClassroomHandler TODO: Add description
func DeleteClassroomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	classroomID := vars["classroom_code"]
	if classroomID == "" {
		http.Error(w, "Missing classroom_code", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.DeleteClassroom(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetClassroomsForTenantHandler TODO: Add description
func GetClassroomsForTenantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Missing tenant_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	classrooms, err := tenantInstance.GetAllClassroomsForTenant()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classrooms)
}

// GetScheduleForTeacherHandler TODO: Add description
func GetScheduleForTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["teacher_id"]
	if teacherID == "" {
		http.Error(w, "Missing teacher_id", http.StatusBadRequest)
		return
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teacher, err := util.GetTeacherByID(teacherID, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	tenants, err := util.GetTenantsForTeacher(teacher, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Failed to get tenants", http.StatusInternalServerError)
		return
	}

	var teacherSchedule []tenantmodels.ScheduleGroup

	for _, tenant := range tenants {
		tenantInstance, err := tenantfactory.Struct(tenant, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		schedule, err := tenantInstance.GetScheduleForTeacher(teacherID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if schedule != nil {
			teacherSchedule = append(teacherSchedule, schedule...)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacherSchedule)
}

// GetTenantsForPupilHandler returns all tenants for a pupil
func GetTenantsForPupilHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Missing pupil_id", http.StatusBadRequest)
		return
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupil, err := util.GetGlobalPupilByID(
		pupilID, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenants, err := util.GetTenantsForPupil(
		*pupil, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}

// GetTenantsForTeacherHandler returns all tenants for a teacher
func GetTenantsForTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["teacher_id"]
	if teacherID == "" {
		http.Error(w, "Missing teacher_id", http.StatusBadRequest)
		return
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teacher, err := util.GetTeacherByID(
		teacherID, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenants, err := util.GetTenantsForTeacher(
		teacher, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenants)
}
