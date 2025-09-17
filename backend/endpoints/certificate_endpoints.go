package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterCertificateEndpoints registers the endpoints for certificate-related operations
func RegisterCertificateEndpoints(r *mux.Router) {
	r.HandleFunc("/api/pupil/certificate/{tenant_id}/{section_id}/{pupil_id}",
		api.AuthMiddleware(
			api.GetCertificateDataHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")
}
