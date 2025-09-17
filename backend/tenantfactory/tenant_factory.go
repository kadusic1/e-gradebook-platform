package tenantfactory

import (
	"database/sql"
	"ednevnik-backend/config"
	wpmodels "ednevnik-backend/models/workspace"
	"ednevnik-backend/tenantshared"
	"ednevnik-backend/util"
	"fmt"
	"net/http"
)

// ConfigurableTenant implements ITenant using configuration
type ConfigurableTenant struct {
	TenantData      wpmodels.Tenant
	Config          config.TenantConfig
	UserWorkspaceDB *sql.DB
	UserTenantDB    *sql.DB
}

// Ensure ConfigurableTenant implements ITenant
var _ tenantshared.ITenant = (*ConfigurableTenant)(nil)

// StructWithDeps TODO: Add description
func StructWithDeps(
	tenant wpmodels.Tenant,
	userAccountType string,
	userWorkspaceDb *sql.DB,
) (tenantshared.ITenant, error) {

	config, exists := config.TenantConfigs[tenant.TenantType]
	if !exists {
		return nil, fmt.Errorf("unsupported tenant type: %s", tenant.TenantType)
	}

	tenantDBName := config.DBPrefix + util.SanitizeString(
		fmt.Sprintf("%d", tenant.ID),
	)

	tenantDB, err := util.GetOrCreateDBConnection(tenantDBName, userAccountType)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant DB: %v", err)
	}

	return &ConfigurableTenant{
		TenantData:      tenant,
		Config:          config,
		UserWorkspaceDB: userWorkspaceDb,
		UserTenantDB:    tenantDB,
	}, nil
}

// Struct keeps old signature, fetches claims and workspace DB
func Struct(tenant wpmodels.Tenant, r *http.Request) (tenantshared.ITenant, error) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing claims in request context")
	}

	userAccountType := claims.AccountType

	userWorkspaceDb, err := util.GetOrCreateDBConnection(
		"ednevnik_workspace", userAccountType,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting workspace DB: %v", err)
	}

	return StructWithDeps(tenant, userAccountType, userWorkspaceDb)
}

// TenantFactory creates appropriate tenant instances
func TenantFactory(tenantID string, r *http.Request) (tenantshared.ITenant, error) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing claims in request context")
	}

	userAccountType := claims.AccountType

	userWorkspaceDb, err := util.GetOrCreateDBConnection(
		"ednevnik_workspace", userAccountType,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting workspace DB: %v", err)
	}

	tenant, err := util.GetTenantByID(tenantID, userWorkspaceDb)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant by ID: %v", err)
	}

	return StructWithDeps(*tenant, userAccountType, userWorkspaceDb)
}

// AccountID creates appropriate tenant instances
func AccountID(tenantID, accountType string) (tenantshared.ITenant, error) {

	userWorkspaceDb, err := util.GetOrCreateDBConnection(
		"ednevnik_workspace", accountType,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting workspace DB: %v", err)
	}

	tenant, err := util.GetTenantByID(tenantID, userWorkspaceDb)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant by ID: %v", err)
	}

	return StructWithDeps(*tenant, accountType, userWorkspaceDb)
}

// CreateDB TODO: Add description
func CreateDB(
	tenant wpmodels.Tenant, r *http.Request,
) (tenantshared.ITenant, error) {
	claims, ok := util.GetClaimsFromContext(r)
	if !ok {
		return nil, fmt.Errorf("unauthorized: missing claims in request context")
	}

	userAccountType := claims.AccountType
	userWorkspaceDb, err := util.GetOrCreateDBConnection(
		"ednevnik_workspace", userAccountType,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting workspace DB: %v", err)
	}

	config, exists := config.TenantConfigs[tenant.TenantType]
	if !exists {
		return nil, fmt.Errorf("unsupported tenant type: %s", tenant.TenantType)
	}

	_, err = util.CreateTenantDB(
		config.DBPrefix,
		fmt.Sprintf("%d", tenant.ID),
		userWorkspaceDb,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating tenant DB: %v", err)
	}

	return StructWithDeps(tenant, userAccountType, userWorkspaceDb)
}

// ServiceReader TODO: Add description
func ServiceReader(tenantID string) (tenantshared.ITenant, error) {
	userWorkspaceDb, err := util.GetOrCreateDBConnectionServiceReader("ednevnik_workspace")
	if err != nil {
		return nil, fmt.Errorf("error getting workspace DB: %v", err)
	}

	tenant, err := util.GetTenantByID(tenantID, userWorkspaceDb)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant by ID: %v", err)
	}

	config, exists := config.TenantConfigs[tenant.TenantType]
	if !exists {
		return nil, fmt.Errorf("unsupported tenant type: %s", tenant.TenantType)
	}

	tenantDBName := config.DBPrefix + util.SanitizeString(
		fmt.Sprintf("%d", tenant.ID),
	)

	tenantDB, err := util.GetOrCreateDBConnectionServiceReader(tenantDBName)
	if err != nil {
		return nil, fmt.Errorf("error getting tenant DB: %v", err)
	}

	return &ConfigurableTenant{
		TenantData:      *tenant,
		Config:          config,
		UserWorkspaceDB: userWorkspaceDb,
		UserTenantDB:    tenantDB,
	}, nil
}
