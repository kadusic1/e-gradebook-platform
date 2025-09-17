package config

// TenantConfig holds configuration for different tenant types
type TenantConfig struct {
	DBPrefix                    string
	SchemaFile                  string
	FinalGradeTable             string
	MaxSemesterCode             string
	AvailableForEnrollmentField string
	BehaviourGradeTable         string
}

// TenantConfigs defines configurations for different tenant types
var TenantConfigs = map[string]TenantConfig{
	"primary": {
		DBPrefix:                    "ednevnik_tenant_db_tenant_id_",
		SchemaFile:                  "db/sql/create_primary_db.sql",
		FinalGradeTable:             "primary_school_final_grades",
		MaxSemesterCode:             "2POL",
		AvailableForEnrollmentField: "available_for_enrollment",
		BehaviourGradeTable:         "primary_school_behaviour_grades",
	},
	"secondary": {
		DBPrefix:                    "ednevnik_tenant_db_tenant_id_",
		SchemaFile:                  "db/sql/create_secondary_db.sql",
		FinalGradeTable:             "high_school_final_grades",
		MaxSemesterCode:             "2POL",
		AvailableForEnrollmentField: "available_for_enrollment",
		BehaviourGradeTable:         "high_school_behaviour_grades",
	},
}
