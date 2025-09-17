package tenantmodels

// PupilSectionInvite TODO: Add description
type PupilSectionInvite struct {
	ID         int    `json:"id"`
	PupilID    int    `json:"pupil_id"`
	SectionID  int    `json:"section_id"`
	InviteDate string `json:"invite_date"`
	Status     string `json:"status"` // ("pending", "accepted", "declined")
	TenantID   int    `json:"tenant_id"`
	// Optional fields
	PupilFullName string `json:"pupil_full_name,omitempty"`
	SectionName   string `json:"section_name,omitempty"`
	TenantName    string `json:"tenant_name,omitempty"`
	PupilEmail    string `json:"pupil_email,omitempty"`
}

// GlobalInvite TODO: Add description
type GlobalInvite struct {
	ID        int `json:"id"`
	InviteID  int `json:"invite_id"`
	AccountID int `json:"account_id"`
	TenantID  int `json:"tenant_id"`
}
