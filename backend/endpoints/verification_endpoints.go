package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterVerificationEndpoints TODO: Add description
func RegisterVerificationEndpoints(r *mux.Router) {
	r.HandleFunc("/api/common/verify_account", api.VerifyAccount).Methods("POST")
}
