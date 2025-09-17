package tenantmodels

// Classroom TODO: Add description
type Classroom struct {
	Code     string `json:"code"`
	Capacity int    `json:"capacity"`
	Type     string `json:"type,omitempty"`
	Name     string `json:"name,omitempty"`
}
