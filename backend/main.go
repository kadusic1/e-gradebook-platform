package main

import (
	"database/sql"
	"ednevnik-backend/api"
	"ednevnik-backend/endpoints"
	"ednevnik-backend/util"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var err error

	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api.JwtKey = []byte(os.Getenv("JWT_SECRET"))

	workspaceCS := util.BuildDBConnectionString("ednevnik_workspace")

	dbWorkspace, err := sql.Open("mysql", workspaceCS)
	if err != nil {
		log.Fatal("Failed to connect to workspace database:", err)
	}
	defer dbWorkspace.Close()

	api.DbWorkspace = dbWorkspace

	r := mux.NewRouter()

	r.Use(api.UserWorkspaceDBMiddleware)

	// CORS setup
	corsAllowedOrigins := handlers.AllowedOrigins([]string{os.Getenv("FRONTEND_URL")})
	corsAllowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})

	r.HandleFunc("/login", api.Login).Methods("POST")
	r.HandleFunc("/parent-login", api.ParentLogin).Methods("POST")

	endpoints.RegisterTeacherEndpoints(r)
	endpoints.RegisterTenantEndpoints(r)
	endpoints.RegisterSectionEndpoints(r)
	endpoints.RegisterStaticEndpoints(r)
	endpoints.RegisterPupilEndpoints(r)
	endpoints.RegisterSemesterEndpoints(r)
	endpoints.RegisterVerificationEndpoints(r)
	endpoints.RegisterDomainEndpoints(r)
	endpoints.RegisterScheduleEndpoints(r)
	endpoints.RegisterClassroomEndpoints(r)
	endpoints.RegisterLessonEndpoints(r)
	endpoints.RegisterGradebookEndpoints(r)
	endpoints.RegisterCertificateEndpoints(r)
	endpoints.RegisterCommonEndpoints(r)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsAllowedOrigins, corsAllowedMethods, corsAllowedHeaders)(r)))
}
