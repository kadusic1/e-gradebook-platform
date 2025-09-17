package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterSectionEndpoints TODO: Add description
func RegisterSectionEndpoints(r *mux.Router) {
	r.HandleFunc("/api/tenant_admin/section_creation_metadata/{tenant_id}",
		api.AuthMiddleware(
			api.GetMetadataForSectionCreation,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/section_create/{tenant_id}",
		api.AuthMiddleware(
			api.CreateSection,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/sections/{tenant_id}/{archived}",
		api.AuthMiddleware(
			api.GetSectionsForTenant,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/section/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.DeleteSection,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/tenant_admin/section/{section_id}",
		api.AuthMiddleware(
			api.UpdateSection,
			[]string{"root", "tenant_admin"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/teacher/send_section_invite/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.SendPupilSectionInvite,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")

	r.HandleFunc("/api/common/respond_to_section_invite",
		api.AuthMiddleware(
			api.PupilInviteHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("POST")

	r.HandleFunc("/api/teacher/unassign_pupil_from_section/{tenant_id}/{section_id}/{pupil_id}",
		api.AuthMiddleware(
			api.UnassignPupilFromSection,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/teacher/archive_section/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.ArchiveSectionHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")

	r.HandleFunc("/api/teacher/unenroll_pupil_from_section/{tenant_id}/{section_id}/{pupil_id}",
		api.AuthMiddleware(
			api.UnenrollPupilFromSectionHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("PUT")
}
