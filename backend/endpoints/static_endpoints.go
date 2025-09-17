package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterStaticEndpoints TODO: Add description
func RegisterStaticEndpoints(r *mux.Router) {
	// Canton endpoint
	r.HandleFunc("/api/common/cantons",
		api.AuthMiddleware(
			api.GetAllCantons,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	// Curriculum endpoints
	r.HandleFunc("/api/tenant_admin/get_curriculums_for_assignment/{tenant_id}",
		api.AuthMiddleware(
			api.GetCurriculumsForTenantAssignment,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/superadmin/assign_curriculums_to_tenant/{tenant_id}",
		api.AuthMiddleware(
			api.AssignCurriculumsToTenant,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/get_curriculums_for_tenant/{tenant_id}",
		api.AuthMiddleware(
			api.GetTenantCurriculums,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/superadmin/unassign_curriculum_from_tenant/{tenant_id}/{curriculum_code}",
		api.AuthMiddleware(
			api.UnassignCurriculumsFromTenant,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/common/get_subjects_for_curriculum/{curriculum_code}",
		api.AuthMiddleware(
			api.GetAllSubjectsForCurriculumHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")
}
