package tenantfactory

import (
	commonmodels "ednevnik-backend/models/common"
	"ednevnik-backend/util"
)

// GetCertificateData retrieves the certificate data for a pupil in a section
func (t *ConfigurableTenant) GetCertificateData(
	sectionID, pupilID int,
) (*commonmodels.Certificate, error) {
	certificate, err := util.GetCertificateData(
		int(t.TenantData.ID),
		sectionID,
		pupilID,
		t.UserTenantDB,
		t.UserWorkspaceDB,
	)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}
