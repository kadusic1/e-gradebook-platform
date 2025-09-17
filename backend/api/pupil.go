package api

import (
	tenantmodels "ednevnik-backend/models/tenant"
	"ednevnik-backend/tenantfactory"
	"ednevnik-backend/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetSectionsForPupilHandler gets all sections the pupil is enrolled in
func GetSectionsForPupilHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	archived := vars["archived"]
	if archived == "" {
		http.Error(w, "Archived parameter is required", http.StatusBadRequest)
		return
	}

	archivedInt, err := strconv.Atoi(archived)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupil, err := util.GetGlobalPupilByID(pupilID, userWorkspaceDb)
	if err != nil {
		http.Error(w, "Pupil not found", http.StatusNotFound)
		return
	}
	if pupil == nil {
		http.Error(w, "Pupil not found", http.StatusNotFound)
		return
	}

	// Get tenants for pupil
	tenants, err := util.GetTenantsForPupil(
		*pupil, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, "Error fetching tenants", http.StatusInternalServerError)
		return
	}

	var pupilSections []tenantmodels.Section

	for _, tenant := range tenants {
		tenantInstance, err := tenantfactory.Struct(tenant, r)
		if err != nil {
			http.Error(w, "Error creating tenant instance", http.StatusInternalServerError)
			return
		}
		sections, err := tenantInstance.GetSectionsForPupil(pupilID, archivedInt)
		if err != nil {
			http.Error(w, "Error fetching sections for pupil", http.StatusInternalServerError)
			return
		}
		for i := range sections {
			sections[i].TenantName = tenant.TenantName
			sections[i].AbsenceDisplay = tenant.AbsenceDisplay
		}
		pupilSections = append(pupilSections, sections...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pupilSections)
}

// GetAbsentAttendancesForPupilHandler retrieves all absent, excused, and unexcused
// pupil attendance records.
func GetAbsentAttendancesForPupilHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Section ID is required", http.StatusBadRequest)
		return
	}

	tenantinstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, "Error creating tenant instance", http.StatusInternalServerError)
		return
	}

	pupilIDint, err := strconv.Atoi(pupilID)
	if err != nil {
		http.Error(w, "Invalid pupil ID", http.StatusBadRequest)
		return
	}
	sectionIDint, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	attendances, err := tenantinstance.GetAbsentAttendancesForPupil(
		pupilIDint, sectionIDint,
	)
	if err != nil {
		http.Error(w, "Error fetching attendances", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendances)
}

// GetSectionBehaviourGradesForPupilHandler used to return behaviour grades
// for a pupil within a section
func GetSectionBehaviourGradesForPupilHandler(w http.ResponseWriter, r *http.Request) {
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

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Missing pupil_id", http.StatusBadRequest)
		return
	}

	pupilIDInt, err := strconv.Atoi(pupilID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	behaviourGrades, err := tenantInstance.GetSectionBehaviourGradesForPupil(
		pupilIDInt, sectionIDInt,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(behaviourGrades)
}

// GetSectionBehaviourGradesForPupilNoPupilIDHandler used to return behaviour grades
// for a pupil within a section without providing pupil ID. We use the pupil ID
// from claims.
func GetSectionBehaviourGradesForPupilNoPupilIDHandler(w http.ResponseWriter, r *http.Request) {
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

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	behaviourGrades, err := tenantInstance.GetSectionBehaviourGradesForPupil(
		claims.ID, sectionIDInt,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(behaviourGrades)
}

// GetCertificateDataHandler retrieves certificate data for a pupil in a section
func GetCertificateDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	sectionID := vars["section_id"]
	if sectionID == "" {
		http.Error(w, "Section ID is required", http.StatusBadRequest)
		return
	}

	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	pupilIDInt, err := strconv.Atoi(pupilID)
	if err != nil {
		http.Error(w, "Invalid pupil ID", http.StatusBadRequest)
		return
	}

	certificate, err := tenantInstance.GetCertificateData(
		sectionIDInt, pupilIDInt,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certificate)
}

// GetGradeEditHistoryHandler retrieves grade edit history
func GetGradeEditHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	gradeID := vars["grade_id"]
	if gradeID == "" {
		http.Error(w, "Grade ID is required", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gradeIDInt, err := strconv.Atoi(gradeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	grades, err := tenantInstance.GetGradeEditHistory(gradeIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// GetBehaviourGradeHistoryHandler retrieves the behaviour grade history for a pupil.
func GetBehaviourGradeHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	behaviourGradeID := vars["behaviour_grade_id"]
	if behaviourGradeID == "" {
		http.Error(w, "Behaviour Grade ID is required", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	behaviourGradeIDInt, err := strconv.Atoi(behaviourGradeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	history, err := tenantInstance.GetBehaviourGradeHistory(behaviourGradeIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func GetLoggedInPupilHandler(w http.ResponseWriter, r *http.Request) {
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

	pupil, err := util.GetGlobalPupilByID(
		fmt.Sprintf("%d", claims.ID), userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, "Pupil not found", http.StatusNotFound)
		return
	}
	if pupil == nil {
		http.Error(w, "Pupil not found", http.StatusNotFound)
		return
	}
	pupil.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pupil)
}

// UpdateLoggedInPupilHandler allows pupils to update their data with service reader connection
func UpdateLoggedInPupilHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var updatedPupil tenantmodels.Pupil
	if err := json.NewDecoder(r.Body).Decode(&updatedPupil); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if claims.ID != updatedPupil.ID {
		http.Error(w, "Unauthorized to update this pupil", http.StatusUnauthorized)
		return
	}

	workspaceDB, err := util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldPupil, err := util.GetGlobalPupilByID(
		fmt.Sprintf("%d", claims.ID), workspaceDB,
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
		tenantInstance, err := tenantfactory.ServiceReader(tenantID)
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

// GetPupilStatisticsFieldsHandler retrieves statistics fields for a pupil
func GetPupilStatisticsFieldsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	pupilIDInt, err := strconv.Atoi(pupilID)
	if err != nil {
		http.Error(w, "Invalid Pupil ID", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Block GET request if claims account type is pupil and pupilIDInt is
	// different from claims ID
	if claims.AccountType == "pupil" && claims.ID != pupilIDInt {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userWorkspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fields, err := util.GetPupilStatisticsFieldsByPupilID(
		pupilIDInt, userWorkspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fields)
}

func UpdatePupilStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	pupilIDInt, err := strconv.Atoi(pupilID)
	if err != nil {
		http.Error(w, "Invalid Pupil ID", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Block PUT request if claims account type is pupil and pupilIDInt is
	// different from claims ID
	if claims.AccountType == "pupil" && claims.ID != pupilIDInt {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var newStatisticsData tenantmodels.PupilStatistics
	if err := json.NewDecoder(r.Body).Decode(&newStatisticsData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	workspaceDB, err := util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedStats, err := util.UpdateStatisticsFieldsForPupil(
		pupilIDInt,
		&newStatisticsData,
		workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStats)
}
