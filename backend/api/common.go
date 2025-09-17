package api

import (
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/tenantfactory"
	"ednevnik-backend/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetAllCantons handles GET /api/cantons and returns all cantons as JSON
func GetAllCantons(w http.ResponseWriter, r *http.Request) {
	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cantons, err := util.GetAllCantons(userWorkspaceDb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cantons)
}

// RegisterTeacher is used to create a new teacher account
func RegisterTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher wpmodels.Teacher
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if teacher.Password == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	err := util.RegisterTeacher(
		teacher,
		DbWorkspace,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// RegisterPupil is used to create a new pupil account
func RegisterPupil(w http.ResponseWriter, r *http.Request) {
	var pupil tenantmodels.Pupil
	if err := json.NewDecoder(r.Body).Decode(&pupil); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if pupil.Password == "" {
		http.Error(w, "Password required", http.StatusBadRequest)
		return
	}

	err := util.RegisterPupil(
		pupil,
		DbWorkspace,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// VerifyAccount is used to verify an account via email
func VerifyAccount(w http.ResponseWriter, r *http.Request) {
	verificationToken := r.URL.Query().Get("token")
	if verificationToken == "" {
		http.Error(w, "Verification token is required", http.StatusBadRequest)
		return
	}

	err := util.VerifyAccount(verificationToken, DbWorkspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetPupilSectionInvites returns all sections invites for a specific pupil
func GetPupilSectionInvites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pupilID, ok := vars["pupil_id"]
	if !ok {
		http.Error(w, "Pupil ID is required", http.StatusBadRequest)
		return
	}

	workspaceDB, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pupil, err := util.GetGlobalPupilByID(pupilID, workspaceDB)
	if err != nil {
		http.Error(w, "Pupil not found", http.StatusNotFound)
		return
	}

	accountID, err := pupil.GetAccountID(workspaceDB)
	if err != nil {
		http.Error(w, "Error fetching account ID", http.StatusInternalServerError)
		return
	}

	globalPupilInvites, err := util.GetGlobalInvitesForAccount(
		accountID,
		workspaceDB,
	)
	if err != nil {
		http.Error(w, "Error fetching global invites", http.StatusInternalServerError)
		return
	}

	var sectionInvites []tenantmodels.PupilSectionInvite

	for _, invite := range globalPupilInvites {
		tenantInstance, err := tenantfactory.ServiceReader(
			fmt.Sprintf("%d", invite.TenantID),
		)
		if err != nil {
			http.Error(w, "Error creating tenant instance", http.StatusInternalServerError)
			return
		}

		invite, err := tenantInstance.GetPupilSectionInvite(
			invite.InviteID,
		)
		if err != nil {
			http.Error(w, "Error fetching pupil section invites", http.StatusInternalServerError)
			return
		}
		sectionInvites = append(sectionInvites, *invite)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sectionInvites)
}

// PupilInviteHandler handles pupil invite accept or decline
func PupilInviteHandler(w http.ResponseWriter, r *http.Request) {
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
		err = tenantInstance.AcceptPupilSectionInvite(
			fmt.Sprintf("%d", inviteRequest.InviteID),
		)
	} else {
		err = tenantInstance.DeclinePupilSectionInvite(
			fmt.Sprintf("%d", inviteRequest.InviteID),
		)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetAllSubjectsForCurriculumHandler returns all subjects for a specific
// curriculum
func GetAllSubjectsForCurriculumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	curriculumCode, ok := vars["curriculum_code"]
	if !ok {
		http.Error(w, "Curriculum code is required", http.StatusBadRequest)
		return
	}

	userWorkspaceDb, ok := util.GetUserWorkspaceDBFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subjects, err := util.GetAllSubjectsForCurriculumCode(
		curriculumCode, userWorkspaceDb,
	)
	if err != nil {
		http.Error(w, "Error fetching subjects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}

func ChangeAccountPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var passwordRequest wpmodels.PasswordChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&passwordRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	workspaceDB, err := util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := util.ChangeAccountPassword(
		claims.AccountID, &passwordRequest, workspaceDB,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ChatHandler handles chatbot requests by forwarding to FastAPI service
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var chatRequest commonmodels.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if chatRequest.Question == "" {
		http.Error(w, "Question is required", http.StatusBadRequest)
		return
	}

	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch tenant names from workspace DB as service reader
	workspaceDB, err := util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
	if err != nil {
		http.Error(w, "Failed to connect to workspace DB", http.StatusInternalServerError)
		return
	}
	var tenantNames []string
	for _, tenantID := range claims.TenantIDs {
		tenant, err := util.GetTenantByID(tenantID, workspaceDB)
		if err == nil && tenant != nil {
			tenantNames = append(tenantNames, tenant.TenantName)
		}
	}

	chatRequest.PermissionData = commonmodels.AIPermissionData{
		ID:                  claims.ID,
		Name:                claims.Name,
		LastName:            claims.LastName,
		Email:               claims.Email,
		Phone:               claims.Phone,
		AccountType:         claims.AccountType,
		AccountID:           claims.AccountID,
		TenantIDs:           claims.TenantIDs,
		TenantAdminTenantID: claims.TenantAdminTenantID,
		TenantNames:         tenantNames,
	}

	// Forward request to FastAPI chatbot service
	chatbotURL := "http://localhost:8005/chat" // Adjust URL as needed

	requestBody, err := json.Marshal(chatRequest)
	if err != nil {
		http.Error(w, "Error preparing request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(chatbotURL, "application/json", strings.NewReader(string(requestBody)))
	if err != nil {
		http.Error(w, "Error communicating with chatbot service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Chatbot service error", resp.StatusCode)
		return
	}

	var chatResponse commonmodels.ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		http.Error(w, "Error parsing chatbot response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse)
}
