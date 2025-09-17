package tenantmodels

// TimePeriod TODO: Add description
type TimePeriod struct {
	ID        int    `json:"id,omitempty"`
	SectionID int    `json:"section_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// Schedule TODO: Add description
type Schedule struct {
	ID            int    `json:"id,omitempty"`
	SectionID     int    `json:"section_id"`
	TimePeriodID  int    `json:"time_period_id,omitempty"`
	SubjectCode   string `json:"subject_code"`
	SubjectName   string `json:"subject_name"`
	Weekday       string `json:"weekday"`
	ClassroomCode string `json:"classroom_code,omitempty"`
	Type          string `json:"type"`
	ColorConfig   int    `json:"color_config,omitempty"`
	TenantName    string `json:"tenant_name,omitempty"`
	SectionName   string `json:"section_name,omitempty"`
	RowStart      string `json:"row_start,omitempty"`
	RowEnd        string `json:"row_end,omitempty"`
}

// ScheduleGroup TODO: Add description
type ScheduleGroup struct {
	TimePeriod TimePeriod `json:"time_period"`
	Schedules  []Schedule `json:"schedules"`
	CreatedAt  string     `json:"created_at,omitempty"`
}

// ScheduleGroupCollection TODO: Add description
type ScheduleGroupCollection []ScheduleGroup
