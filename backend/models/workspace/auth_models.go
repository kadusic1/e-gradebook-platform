package wpmodels

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims for a teacher
// All teacher data is included in the token
// (Password is omitted for security)
type Claims struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	LastName            string   `json:"last_name"`
	Email               string   `json:"email"`
	Phone               string   `json:"phone"`
	AccountType         string   `json:"account_type"`
	AccountID           int      `json:"account_id"`
	TenantIDs           []string `json:"tenant_ids"`
	TenantAdminTenantID int      `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

// AuthRequest - DRY concept: Define a struct for the authentication request
// Useful for register and login endpoints
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ContextKey TODO: Add description
type ContextKey string

type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
