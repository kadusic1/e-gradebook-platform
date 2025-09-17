package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterClassroomEndpoints TODO: Add description
func RegisterClassroomEndpoints(r *mux.Router) {
	r.HandleFunc("/api/tenant_admin/create_classroom/{tenant_id}",
		api.AuthMiddleware(
			api.CreateClassroomHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/tenant_admin/update_classroom/{tenant_id}/{classroom_code}",
		api.AuthMiddleware(
			api.UpdateClassroomHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/tenant_admin/delete_classroom/{tenant_id}/{classroom_code}",
		api.AuthMiddleware(
			api.DeleteClassroomHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/tenant_admin/update_classroom/{tenant_id}/{classroom_code}",
		api.AuthMiddleware(
			api.DeleteClassroomHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/tenant_admin/get_all_classrooms/{tenant_id}",
		api.AuthMiddleware(
			api.GetClassroomsForTenantHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")
}
