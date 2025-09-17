package util

import (
	"database/sql"
	"fmt"
	"net/http"

	"ednevnik-backend/constants"
	"ednevnik-backend/models/interfaces"
	wpmodels "ednevnik-backend/models/workspace"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ParseAndValidateJWT parses the JWT token string and returns the claims if valid, otherwise an error.
func ParseAndValidateJWT(tokenStr string, jwtKey []byte) (*wpmodels.Claims, error) {
	// Strip "Bearer " prefix if present
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	token, err := jwt.ParseWithClaims(tokenStr, &wpmodels.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(*wpmodels.Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}

// ComparePassword compares a bcrypt hashed password with its possible plaintext equivalent.
func ComparePassword(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// GetClaimsFromContext extracts all JWT claims from the request context.
func GetClaimsFromContext(r *http.Request) (*wpmodels.Claims, bool) {
	claims, ok := r.Context().Value(constants.ClaimsKey).(*wpmodels.Claims)
	return claims, ok && claims != nil
}

// GetUserByEmail retrieves a user by email, checking teachers first, then pupils
func GetUserByEmail(email string, workspaceDB *sql.DB) (interfaces.User, error) {
	// Try to get user as teacher first
	if user, err := GetTeacherByEmail(workspaceDB, email); err == nil && user != nil {
		return user, nil
	}

	// Try to get user as pupil
	if user, err := GetGlobalPupilByEmail(email, workspaceDB); err == nil && user != nil {
		return user, nil
	}

	return nil, sql.ErrNoRows
}

// ChangeAccountPassword changes the password for a specific account.
func ChangeAccountPassword(
	accountID int,
	passwordRequest *wpmodels.PasswordChangeRequest,
	workspaceDB *sql.DB,
) error {
	if passwordRequest.NewPassword != passwordRequest.ConfirmPassword {
		return fmt.Errorf("nove lozinke se ne podudaraju. Molimo pokušajte ponovo")
	}

	var currentHashedPassword string
	err := workspaceDB.QueryRow(
		"SELECT password FROM accounts WHERE id = ?",
		accountID,
	).Scan(&currentHashedPassword)
	if err != nil {
		return err
	}

	// Compare current password
	if err := ComparePassword(currentHashedPassword, passwordRequest.CurrentPassword); err != nil {
		return fmt.Errorf("unesena trenutna lozinka nije ispravna")
	}

	// Hash new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(passwordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password in DB
	_, err = workspaceDB.Exec(
		"UPDATE accounts SET password = ? WHERE id = ?",
		string(hashedNewPassword),
		accountID,
	)
	if err != nil {
		return err
	}

	return nil
}
