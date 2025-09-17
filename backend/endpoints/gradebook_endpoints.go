package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterGradebookEndpoints function to register endpoints related to the gradebook
func RegisterGradebookEndpoints(r *mux.Router) {
	r.HandleFunc("/api/teacher/gradebook_metadata/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GradebookMetadataHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/section_grades_for_subject/{tenant_id}/{section_id}/{subject_code}/{semester_code}",
		api.AuthMiddleware(
			api.GetSectionGradesForSubjectHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/create_grade/{tenant_id}",
		api.AuthMiddleware(
			api.CreateGradeHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("POST")

	r.HandleFunc("/api/teacher/delete_grade/{tenant_id}",
		api.AuthMiddleware(
			api.DeleteGradeHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("DELETE")

	r.HandleFunc("/api/teacher/update_grade/{tenant_id}",
		api.AuthMiddleware(
			api.UpdateGradeHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("PUT")

	r.HandleFunc("/api/pupil/section_grades_for_pupil/{tenant_id}/{section_id}/{semester_code}",
		api.AuthMiddleware(
			api.GetPupilGradesForSectionPupilHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/behaviour_grades/{tenant_id}/{section_id}/{pupil_id}",
		api.AuthMiddleware(
			api.GetSectionBehaviourGradesForPupilHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/behaviour_grades/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetSectionBehaviourGradesForPupilNoPupilIDHandler,
			[]string{"pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/grade_edit_history/{tenant_id}/{grade_id}",
		api.AuthMiddleware(
			api.GetGradeEditHistoryHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/pupil/behaviour_grade_history/{tenant_id}/{behaviour_grade_id}",
		api.AuthMiddleware(
			api.GetBehaviourGradeHistoryHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("GET")

	r.HandleFunc("/api/teacher/complete_gradebook/{tenant_id}/{section_id}",
		api.AuthMiddleware(
			api.GetCompleteGradebookDataHandler,
			[]string{"root", "tenant_admin", "teacher"},
		),
	).Methods("GET")
}
