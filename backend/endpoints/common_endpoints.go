package endpoints

import (
	"ednevnik-backend/api"

	"github.com/gorilla/mux"
)

// RegisterCommonEndpoints registers the endpoints for common operations
func RegisterCommonEndpoints(r *mux.Router) {
	r.HandleFunc("/api/common/change_password",
		api.AuthMiddleware(
			api.ChangeAccountPasswordHandler,
			[]string{"teacher", "pupil"},
		),
	).Methods("POST")

	r.HandleFunc("/api/common/chat",
		api.AuthMiddleware(
			api.ChatHandler,
			[]string{"root", "tenant_admin", "teacher", "pupil"},
		),
	).Methods("POST")
}
