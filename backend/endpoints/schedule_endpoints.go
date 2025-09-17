package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterScheduleEndpoints TODO: Add description
func RegisterScheduleEndpoints(r *mux.Router) {
	r.HandleFunc("/api/tenant_admin/schedule_create/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.CreateScheduleHandler,
			[]string{"root", "tenant_admin"},
		),
	).Methods("POST")

	r.HandleFunc("/api/pupil/schedule/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetScheduleForSectionHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/schedule/{teacher_id}",
		api.AuthMiddleware(
			api.GetScheduleForTeacherHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")
}
