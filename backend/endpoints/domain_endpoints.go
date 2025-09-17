package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterDomainEndpoints TODO: Add description
func RegisterDomainEndpoints(r *mux.Router) {
	r.HandleFunc("/api/superadmin/all_domains",
		api.AuthMiddleware(
			api.GetAllDomainsHandler,
			[]string{"root"},
		),
	).Methods("GET")

	r.HandleFunc("/api/superadmin/domain_create",
		api.AuthMiddleware(
			api.CreateGlobalDomainHandler,
			[]string{"root"},
		),
	).Methods("POST")

	r.HandleFunc("/api/superadmin/domains_delete",
		api.AuthMiddleware(
			api.DeleteGlobalDomainHandler,
			[]string{"root"},
		),
	).Methods("DELETE")
}
