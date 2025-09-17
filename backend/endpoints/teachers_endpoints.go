package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterTeacherEndpoints TODO: Add description
func RegisterTeacherEndpoints(r *mux.Router) {
	// Teacher CRUD routes
	r.HandleFunc("/api/tenant_admin/teacher",
		api.AuthMiddleware(
			api.CreateTeacher,
			[]string{"tenant_admin", "root"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/teachers",
		api.AuthMiddleware(
			api.ListTeachers,
			[]string{"tenant_admin", "root", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/get_teacher/{id}",
		api.AuthMiddleware(
			api.GetTeacher,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/teacher/{id}",
		api.AuthMiddleware(
			api.UpdateTeacher,
			[]string{"tenant_admin", "root", "teacher"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/tenant_admin/teacher/{id}",
		api.AuthMiddleware(
			api.DeleteTeacher,
			[]string{"tenant_admin", "root"},
		),
	).Methods("DELETE")

	// Teacher assignment routes
	r.HandleFunc("/api/tenant_admin/teachers_per_tenant/{tenant_id}",
		api.AuthMiddleware(
			api.ListTeachersPerTenant,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	// Teacher registration route
	r.HandleFunc("/api/common/register_teacher", api.RegisterTeacher).Methods("POST")

	// Teacher section assignment routes
	r.HandleFunc("/api/tenant_admin/teacher_section_assignments/{tenant_id}/{teacher_id}",
		api.AuthMiddleware(
			api.HandleTeacherSectionAssignments,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/teacher/invites/{teacher_id}",
		api.AuthMiddleware(
			api.GetTeacherInvitesHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/teacher_invites/{tenant_id}",
		api.AuthMiddleware(
			api.GetTeacherInvitesForTenantHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/handle_invite",
		api.AuthMiddleware(
			api.TeacherInviteHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/delete_teacher_invite",
		api.AuthMiddleware(
			api.DeleteTeacherInviteHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/tenant_admin/delete_teacher_from_tenant/{tenant_id}/{teacher_id}",
		api.AuthMiddleware(
			api.DeleteTeacherForTenantHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/teacher/sections/{teacher_id}/{archived}",
		api.AuthMiddleware(
			api.GetSectionsForTeacherHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/tenant_admin/tenants_for_teacher/{teacher_id}",
		api.AuthMiddleware(
			api.GetTenantsForTeacherHandler, []string{"root", "tenant_admin"},
		),
	).Methods("GET")
}
