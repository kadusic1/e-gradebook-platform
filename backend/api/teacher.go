package api

import (
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

// GetSectionsForTeacherHandler TODO: Add description
func GetSectionsForTeacherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["teacher_id"]
	if teacherID == "" {
		http.Error(w, "Teacher ID is required", http.StatusBadRequest)
		return
	}

	archived := vars["archived"]
	if archived == "" {
		http.Error(w, "Archived parameter is required", http.StatusBadRequest)
		return
	}

	archivedInt, err := strconv.Atoi(archived)
	if err != nil {
		http.Error(w, "Error converting archived parameter type", http.StatusNotFound)
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

	// Get tenants for teacher
	tenants, err := util.GetTenantsForTeacher(
		teacher, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, "Error fetching tenants", http.StatusInternalServerError)
		return
	}

	var teacherSections []tenantmodels.Section

	for _, tenant := range tenants {
		tenantInstance, err := tenantfactory.Struct(tenant, r)
		if err != nil {
			http.Error(w, "Error creating tenant instance", http.StatusInternalServerError)
			return
		}
		sections, err := tenantInstance.GetSectionsForTeacher(teacherID, archivedInt)
		if err != nil {
			http.Error(w, "Error fetching sections for teacher", http.StatusInternalServerError)
			return
		}
		for i := range sections {
			sections[i].TenantName = tenant.TenantName
			sections[i].PupilDisplay = tenant.PupilDisplay
			sections[i].PupilInviteDisplay = tenant.PupilInviteDisplay
			sections[i].LessonDisplay = tenant.LessonDisplay
			sections[i].AbsenceDisplay = tenant.AbsenceDisplay
		}
		teacherSections = append(teacherSections, sections...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacherSections)
}

// GetPupilsForSection TODO: Add description
func GetPupilsForSection(w http.ResponseWriter, r *http.Request) {
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

	pupils, err := tenantInstance.GetPupilsForSection(sectionID, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*pupils)
}

// DeletePupilInviteHandler TODO: Add description
func DeletePupilInviteHandler(w http.ResponseWriter, r *http.Request) {
	type DeleteRequest struct {
		InviteID int `json:"invite_id"`
		PupilID  int `json:"pupil_id"`
		TenantID int `json:"tenant_id"`
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

	err = tenantInstance.DeletePupilInvite(
		request.InviteID,
		request.PupilID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupil, err := util.GetGlobalPupilByID(
		fmt.Sprintf("%d", request.PupilID), workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pupil.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pupil)
}

// UnassignPupilFromSection TODO: Add description
func UnassignPupilFromSection(w http.ResponseWriter, r *http.Request) {
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

	pupilID := vars["pupil_id"]
	if pupilID == "" {
		http.Error(w, "Missing pupil_id", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.DeletePupilFromSection(pupilID, sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTeacherInvitesHandler TODO: Add description
func GetTeacherInvitesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teacherID := vars["teacher_id"]
	if teacherID == "" {
		http.Error(w, "Missing teacher_id", http.StatusBadRequest)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the teacher object
	teacher, err := util.GetTeacherByID(teacherID, workspaceDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get teacher account id
	teacherAccountID, err := teacher.GetAccountID(workspaceDB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get global invites for teacher
	globalInvites, err := util.GetGlobalInvitesForAccount(
		teacherAccountID, workspaceDB,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var teacherInvites []wpmodels.TeacherSectionInviteRecord

	for _, invite := range globalInvites {
		tenantInstance, err := tenantfactory.ServiceReader(
			fmt.Sprintf("%d", invite.TenantID),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tenantInvites, err := tenantInstance.GetInvitesForTeacher(
			invite.InviteID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teacherInvites = append(teacherInvites, tenantInvites...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacherInvites)
}

// TeacherInviteHandler TODO: Add description
func TeacherInviteHandler(w http.ResponseWriter, r *http.Request) {
	type InviteRequest struct {
		InviteID int    `json:"invite_id"`
		Action   string `json:"action"` // "accept" or "decline"
		TenantID int    `json:"tenant_id"`
	}

	var inviteRequest InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&inviteRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.ServiceReader(
		fmt.Sprintf("%d", inviteRequest.TenantID),
	)
	if err != nil {
		http.Error(w, "Error creating tenant instance", http.StatusInternalServerError)
		return
	}

	if inviteRequest.Action == "accept" {
		err = tenantInstance.AcceptTeacherSectionInvite(
			fmt.Sprintf("%d", inviteRequest.InviteID),
		)
	} else {
		err = tenantInstance.DeclineTeacherSectionInvite(
			fmt.Sprintf("%d", inviteRequest.InviteID),
		)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetLessonsForSectionHandler retrieves all lessons for a specific section within a tenant.
// Expects tenant_id and section_id as URL parameters.
// Returns JSON array of lesson data on success.
func GetLessonsForSectionHandler(w http.ResponseWriter, r *http.Request) {
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

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	lessons, err := tenantInstance.GetLessonsForSection(sectionIDInt, claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pupils, err := tenantInstance.GetPupilsForSection(sectionID, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims, ok = util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subjects, err := tenantInstance.GetSubjectsForSection(
		sectionIDInt, claims,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type lessonResponse struct {
		Pupils   []tenantmodels.Pupil      `json:"pupils"`
		Lessons  []tenantmodels.LessonData `json:"lessons"`
		Subjects []wpmodels.Subject        `json:"subjects"`
	}

	response := lessonResponse{
		Pupils:   pupils.Pupils,
		Lessons:  lessons,
		Subjects: subjects,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateSectionLessonHandler creates a new lesson for a specific section within a tenant.
// Expects tenant_id and section_id as URL parameters and lesson data in JSON request body.
// Returns HTTP 201 Created status on successful creation.
func CreateSectionLessonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var lessonData tenantmodels.LessonData
	err := json.NewDecoder(r.Body).Decode(&lessonData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newLesson, err := tenantInstance.CreateSectionLesson(lessonData, claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newLesson)
}

// UpdateLessonHandler updates an existing lesson by ID within a tenant.
// Expects tenant_id and lesson_id as URL parameters and updated lesson data in JSON request body.
// Returns HTTP 204 No Content status on successful update.
func UpdateLessonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	lessonID := vars["lesson_id"]
	if lessonID == "" {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var lessonData tenantmodels.LessonData
	err := json.NewDecoder(r.Body).Decode(&lessonData)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lessonIDInt, err := strconv.Atoi(lessonID)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}

	updatedLesson, err := tenantInstance.UpdateLesson(lessonIDInt, lessonData, claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedLesson)
}

// DeleteLessonHandler deletes a lesson by ID within a tenant.
// Expects tenant_id and lesson_id as URL parameters.
// Returns HTTP 204 No Content status on successful deletion.
func DeleteLessonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	lessonID := vars["lesson_id"]
	if lessonID == "" {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lessonIDInt, err := strconv.Atoi(lessonID)
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}

	err = tenantInstance.DeleteLesson(lessonIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAbsentAttendancesForSectionHandler retrieves all absent, excused, and unexcused
// pupil attendance records for a specific section within a tenant.
func GetAbsentAttendancesForSectionHandler(w http.ResponseWriter, r *http.Request) {
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

	attendances, err := tenantInstance.GetAbsentAttendancesForSection(sectionIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attendances)
}

// HandleAttendanceActionHandler processes attendance actions
func HandleAttendanceActionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	var action tenantmodels.AttendanceAction
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.HandleAttendanceAction(action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GradebookMetadataHandler retrieves metadata for the gradebook,
// including subjects and pupil count.
func GradebookMetadataHandler(w http.ResponseWriter, r *http.Request) {
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

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	var subjects []wpmodels.Subject
	var pupilCount int

	if claims.AccountType != "pupil" {

		subjects, err = tenantInstance.GetSubjectsForSection(
			sectionIDInt, claims,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pupilCount, err = tenantInstance.GetPupilCountForSection(sectionIDInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	sectionSemesters, err := tenantInstance.GetSemestersForSection(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type responseData struct {
		Subjects         []wpmodels.Subject        `json:"subjects,omitempty"`
		PupilCount       int                       `json:"pupil_count"`
		SectionSemesters []wpmodels.TenantSemester `json:"section_semesters,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	if claims.AccountType != "pupil" {
		json.NewEncoder(w).Encode(responseData{
			Subjects:         subjects,
			PupilCount:       pupilCount,
			SectionSemesters: sectionSemesters,
		})
	} else {
		json.NewEncoder(w).Encode(responseData{
			SectionSemesters: sectionSemesters,
		})
	}
}

// GetSectionGradesForSubjectHandler retrieves grades for all pupils in a section
// for a specific subject.
func GetSectionGradesForSubjectHandler(w http.ResponseWriter, r *http.Request) {
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

	subjectCode := vars["subject_code"]
	if subjectCode == "" {
		http.Error(w, "Subject code is required", http.StatusBadRequest)
		return
	}

	semesterCode := vars["semester_code"]
	if semesterCode == "" || semesterCode == "undefined" {
		http.Error(w, "Semester code is required", http.StatusBadRequest)
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

	grades, err := tenantInstance.GetSectionGradesForSubject(
		sectionIDInt, semesterCode, subjectCode,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// CreateGradeHandler creates a new grade for a pupil in a section and subject.
func CreateGradeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	var grade tenantmodels.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	grade.TeacherID = claims.ID

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdGrade, err := tenantInstance.CreateGrade(&grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdGrade)
}

// DeleteGradeHandler deletes a grade for a pupil in a section and subject.
func DeleteGradeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var grade tenantmodels.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gradesAfterDeletion, err := tenantInstance.DeleteGrade(&grade, claims.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gradesAfterDeletion)
}

// UpdateGradeHandler updates a grade
func UpdateGradeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	var grade tenantmodels.Grade
	if err := json.NewDecoder(r.Body).Decode(&grade); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	grade.TeacherID = claims.ID

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedGradeItems, err := tenantInstance.UpdateGrade(&grade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedGradeItems)
}

// GetPupilGradesForSectionPupilHandler retrieves pupil grades for all subjects
// in a section
func GetPupilGradesForSectionPupilHandler(w http.ResponseWriter, r *http.Request) {
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

	semesterCode := vars["semester_code"]
	if semesterCode == "" || semesterCode == "undefined" {
		http.Error(w, "Semester code is required", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	grades, err := tenantInstance.GetPupilGradesForSectionPupil(
		sectionIDInt, claims.ID, semesterCode,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grades)
}

// UpdatePupilBehaviourGradeHandler updates a pupil behaviour grade
func UpdatePupilBehaviourGradeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	if tenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var behaviourGradesToupdate tenantmodels.BehaviourGrade
	if err := json.NewDecoder(r.Body).Decode(&behaviourGradesToupdate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tenantInstance, err := tenantfactory.TenantFactory(
		tenantID, r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedBehaviourGrade, err := tenantInstance.UpdatePupilBehaviourGrade(
		behaviourGradesToupdate, claims.ID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBehaviourGrade)
}

// ArchiveSectionHandler archives a section with a specific section ID
func ArchiveSectionHandler(w http.ResponseWriter, r *http.Request) {
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

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tenantInstance, err := tenantfactory.TenantFactory(
		tenantID, r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.ArchiveSection(sectionIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCompleteGradebookDataHandler returns complete gradebook data for a section.
func GetCompleteGradebookDataHandler(w http.ResponseWriter, r *http.Request) {
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

	tenantInstance, err := tenantfactory.TenantFactory(
		tenantID, r,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		http.Error(w, "Invalid section ID", http.StatusBadRequest)
		return
	}

	completeGradebook, err := tenantInstance.GetCompleteGradebookData(sectionIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completeGradebook)
}

// UnenrollPupilFromSectionHandler unenrolls a pupil from a section
func UnenrollPupilFromSectionHandler(w http.ResponseWriter, r *http.Request) {
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

	tenantInstance, err := tenantfactory.TenantFactory(tenantID, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tenantInstance.UnenrollPupilFromSection(pupilIDInt, sectionIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
