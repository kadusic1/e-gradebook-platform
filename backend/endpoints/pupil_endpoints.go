package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterPupilEndpoints TODO: Add description
func RegisterPupilEndpoints(r *mux.Router) {
	r.HandleFunc("/api/teacher/pupils/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetPupilsForSection, []string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/pupils",
		api.AuthMiddleware(
			api.CreatePupil,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/pupils/{pupil_id}",
		api.AuthMiddleware(
			api.DeletePupil,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/tenant_admin/pupils/{pupil_id}",
		api.AuthMiddleware(
			api.UpdatePupil,
			[]string{"root", "tenant_admin"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/tenant_admin/pupil_accounts",
		api.AuthMiddleware(
			api.ListPupilAccounts,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/common/register_pupil", api.RegisterPupil).Methods("POST")

	r.HandleFunc("/api/common/pupil_section_invites/{pupil_id}",
		api.AuthMiddleware(
			api.GetPupilSectionInvites,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/delete_pupil_invite",
		api.AuthMiddleware(
			api.DeletePupilInviteHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/pupil/sections/{pupil_id}/{archived}",
		api.AuthMiddleware(
			api.GetSectionsForPupilHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/tenants_for_pupil/{pupil_id}",
		api.AuthMiddleware(
			api.GetTenantsForPupilHandler, []string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/update_behaviour_grade/{tenant_id}",
		api.AuthMiddleware(
			api.UpdatePupilBehaviourGradeHandler, []string{"root", "tenant_admin", "teacher"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/pupil/profile",
		api.AuthMiddleware(
			api.GetLoggedInPupilHandler, []string{"pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/update_general_data",
		api.AuthMiddleware(
			api.UpdateLoggedInPupilHandler, []string{"pupil"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/pupil/statistics_fields/{pupil_id}",
		api.AuthMiddleware(
			api.GetPupilStatisticsFieldsHandler, []string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/update_statistics_fields/{pupil_id}",
		api.AuthMiddleware(
			api.UpdatePupilStatisticsHandler, []string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("PUT")
}
