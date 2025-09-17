package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ednevnik-backend/constants"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/util"

	"github.com/golang-jwt/jwt/v5"
)

// JwtKey is the secret key used to sign JWT tokens
// It will be read from an environment variable in the main package
var JwtKey []byte

// DbWorkspace is the connection to the workspace database
var DbWorkspace *sql.DB

// Login allows user login
func Login(w http.ResponseWriter, r *http.Request) {
	var req wpmodels.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	email := req.Email
	password := req.Password
	if email == "" || password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		return
	}

	domains, err := util.GetAllDomainsHelper(DbWorkspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the email ends with a valid domain
	validDomain := false
	for _, domain := range domains {
		if strings.HasSuffix(req.Email, domain.Domain) {
			validDomain = true
			break
		}
	}
	if !validDomain {
		http.Error(w, "Invalid domain", http.StatusBadGateway)
		return
	}

	user, err := util.GetUserByEmail(email, DbWorkspace)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare hashed password using bcrypt
	if err := util.ComparePassword(user.GetPassword(), password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userTenantIDs, err := user.GetTenantIDs(DbWorkspace)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	userAccountID, err := user.GetAccountID(DbWorkspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accountType := user.GetAccountType(DbWorkspace)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &wpmodels.Claims{
		ID:          user.GetID(),
		Name:        user.GetName(),
		LastName:    user.GetLastName(),
		Email:       user.GetEmail(),
		Phone:       user.GetPhone(),
		AccountType: accountType,
		AccountID:   userAccountID,
		TenantIDs:   userTenantIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	if accountType == "tenant_admin" {
		tenantID, err := util.GetTenantIDForTenantAdmin(
			user.GetID(), DbWorkspace,
		)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		claims.TenantAdminTenantID = tenantID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

// ParentLogin allows parent login using access code
func ParentLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ParentAccessCode string `json:"parent_access_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	parentAccessCode := req.ParentAccessCode
	if parentAccessCode == "" {
		http.Error(w, "Missing parent access code", http.StatusBadRequest)
		return
	}

	// Get pupil by parent access code
	pupil, err := util.GetGlobalPupilByParentAccessCode(parentAccessCode, DbWorkspace)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if pupil == nil {
		http.Error(w, "Invalid parent access code", http.StatusUnauthorized)
		return
	}

	pupilTenantIDs, err := pupil.GetTenantIDs(DbWorkspace)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	pupilAccountID, err := pupil.GetAccountID(DbWorkspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accountType := pupil.GetAccountType(DbWorkspace)

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &wpmodels.Claims{
		ID:          pupil.ID,
		Name:        pupil.Name,
		LastName:    pupil.LastName,
		Email:       pupil.Email,
		Phone:       pupil.PhoneNumber,
		AccountType: accountType,
		AccountID:   pupilAccountID,
		TenantIDs:   pupilTenantIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

// AuthMiddleware is used for user authentication, one of roles:
// root, tenant_admin, teacher, pupil
func AuthMiddleware(next http.HandlerFunc, accountTypes []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		claims, err := util.ParseAndValidateJWT(tokenStr, JwtKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		allowed := false
		for _, allowedType := range accountTypes {
			if strings.EqualFold(claims.AccountType, allowedType) {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constants.ClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// UserWorkspaceDBMiddleware is used to get a workspaceDB instance
// for the logged in user
func UserWorkspaceDBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" || r.URL.Path == "/parent-login" {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}
		claims, err := util.ParseAndValidateJWT(tokenStr, JwtKey)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		userDB, err := util.GetOrCreateDBConnection("ednevnik_workspace", claims.AccountType)
		if err != nil {
			http.Error(w, "Failed to connect as user", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserWorkspaceDBKey, userDB)
		ctx = context.WithValue(ctx, constants.ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
