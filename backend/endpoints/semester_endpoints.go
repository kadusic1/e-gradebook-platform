package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterSemesterEndpoints TODO: Add description
func RegisterSemesterEndpoints(r *mux.Router) {
	r.HandleFunc("/api/superadmin/npp_semesters",
		api.AuthMiddleware(
			api.GetAllNPPSemesters,
			[]string{"root"},
		),
	).Methods("GET")

	r.HandleFunc("/api/superadmin/npp_semesters_update",
		api.AuthMiddleware(
			api.UpdateNPPSemesterDates,
			[]string{"root"},
		),
	).Methods("PUT")
}
