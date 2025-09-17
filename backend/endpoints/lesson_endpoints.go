package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterLessonEndpoints TODO: Add description
func RegisterLessonEndpoints(r *mux.Router) {
	r.HandleFunc("/api/teacher/get_lessons_for_section/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetLessonsForSectionHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/create_lesson/{tenant_id}",
		api.AuthMiddleware(
			api.CreateSectionLessonHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")

	r.HandleFunc("/api/teacher/update_lesson/{tenant_id}/{lesson_id}",
		api.AuthMiddleware(
			api.UpdateLessonHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/teacher/delete_lesson/{tenant_id}/{lesson_id}",
		api.AuthMiddleware(
			api.DeleteLessonHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/teacher/get_absent_attendances_for_section/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetAbsentAttendancesForSectionHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/get_absent_attendances_for_pupil/{tenant_id}/{section_id}/{pupil_id}",
		api.AuthMiddleware(
			api.GetAbsentAttendancesForPupilHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/handle_attendance_action/{tenant_id}",
		api.AuthMiddleware(
			api.HandleAttendanceActionHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")
}
