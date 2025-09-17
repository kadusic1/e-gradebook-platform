package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterTenantEndpoints TODO: Add description
func RegisterTenantEndpoints(r *mux.Router) {
	// Tenant CRUD routes
	r.HandleFunc("/api/superadmin/tenant",
		api.AuthMiddleware(
			api.CreateTenant,
			[]string{"root"},
		),
	).Methods("POST")

	r.HandleFunc("/api/superadmin/tenants",
		api.AuthMiddleware(
			api.ListTenants,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/tenant/{id}",
		api.AuthMiddleware(
			api.GetTenant,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/superadmin/tenant/{id}",
		api.AuthMiddleware(
			api.UpdateTenant,
			[]string{"root"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/superadmin/tenant/{id}",
		api.AuthMiddleware(
			api.DeleteTenant,
			[]string{"root"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/tenant_admin/tenant_semesters/{tenant_id}",
		api.AuthMiddleware(
			api.GetTenantSemesters,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/tenant_semesters_update",
		api.AuthMiddleware(
			api.UpdateTenantSemesterDates,
			[]string{"root", "tenant_admin"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/tenant_admin/teacher_invite_data/{tenant_id}",
		api.AuthMiddleware(
			api.GetTeacherInviteDataForTenantHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")
}
