package util

import (
	"database/sql"
	tenantmodels "ednevnik-backend/models/tenant"

	"github.com/google/uuid"
)

// CreateTimePeriod TODO: Add description
func CreateTimePeriod(
	timePeriod tenantmodels.TimePeriod, tenantDB *sql.Tx, batchID string,
) (int, error) {
	query := `INSERT INTO time_periods (section_id, start_time, end_time, batch_id)
	VALUES (?, ?, ?, ?)`
	result, err := tenantDB.Exec(
		query, timePeriod.SectionID, timePeriod.StartTime, timePeriod.EndTime, batchID,
	)
	if err != nil {
		return 0, err
	}
	timePeriodID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(timePeriodID), nil
}

var weekdayConvertToEnglishMap = map[string]string{
	"ponedjeljak": "Monday",
	"utorak":      "Tuesday",
	"srijeda":     "Wednesday",
	"četvrtak":    "Thursday",
	"petak":       "Friday",
}

// CreateScheduleItem TODO: Add description
func CreateScheduleItem(
	item tenantmodels.Schedule, tenantDB *sql.Tx, batchID string,
) error {
	item.Weekday = weekdayConvertToEnglishMap[item.Weekday]

	var classroomCode interface{}
	if item.ClassroomCode == "" || item.ClassroomCode == "null" {
		classroomCode = nil
	} else {
		classroomCode = item.ClassroomCode
	}

	query := `INSERT INTO schedule (section_id, time_period_id, subject_code,
	weekday, classroom_code, type, batch_id)
	VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := tenantDB.Exec(
		query, item.SectionID, item.TimePeriodID, item.SubjectCode,
		item.Weekday, classroomCode, "regular", batchID,
	)
	return err
}

// CreateSchedule TODO: Add description
func CreateSchedule(
	scheduleData tenantmodels.ScheduleGroupCollection,
	tenantDB *sql.DB,
	sectionID string,
) (err error) {
	// Start a transaction
	tx, err := tenantDB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	batchID := uuid.NewString()

	// If a schedule already exists for the section, delete it
	err = DeleteSchedule(sectionID, tenantDB)
	if err != nil {
		return err
	}

	for _, timePeriodGroup := range scheduleData {
		// First insert the time period
		timePeriodInt, err := CreateTimePeriod(
			timePeriodGroup.TimePeriod, tx, batchID,
		)
		if err != nil {
			return err
		}

		for _, scheduleItem := range timePeriodGroup.Schedules {
			// Set the time period ID for the schedule item
			scheduleItem.TimePeriodID = timePeriodInt

			// Then insert the schedule item
			err = CreateScheduleItem(scheduleItem, tx, batchID)
			if err != nil {
				return err
			}
		}

	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteSchedule TODO: Add description
func DeleteSchedule(
	sectionID string, tenantDB *sql.DB,
) error {
	query := `DELETE FROM time_periods WHERE section_id = ?`
	_, err := tenantDB.Exec(query, sectionID)
	return err
}

var weekdayConvertToBosnianMap = map[string]string{
	"Monday":    "ponedjeljak",
	"Tuesday":   "utorak",
	"Wednesday": "srijeda",
	"Thursday":  "četvrtak",
	"Friday":    "petak",
}

// processScheduleRows processes the result rows from a schedule query and returns grouped schedule data
func processScheduleRows(rows *sql.Rows) (tenantmodels.ScheduleGroupCollection, error) {
	defer rows.Close()

	// Map to group schedules by time period ID
	timePeriodMap := make(map[int]*tenantmodels.ScheduleGroup)
	var timePeriodOrder []int // To maintain order

	for rows.Next() {
		var tpID, tpSectionID int
		var tpStartTime, tpEndTime string
		var sID, sSectionID, sTimePeriodID, tColorConfig sql.NullInt64
		var sSubjectCode, sSubjectName, sWeekday, sClassroomCode, sType sql.NullString
		var tenantName, sectionName sql.NullString

		err := rows.Scan(
			&tpID, &tpSectionID, &tpStartTime, &tpEndTime,
			&sID, &sSectionID, &sTimePeriodID, &sSubjectCode, &sSubjectName,
			&sWeekday, &sClassroomCode, &sType, &tColorConfig, &tenantName, &sectionName,
		)
		if err != nil {
			return nil, err
		}

		// Convert weekday to Bosnian
		if sWeekday.Valid {
			sWeekday.String = weekdayConvertToBosnianMap[sWeekday.String]
		}

		// Check if we already have this time period
		if _, exists := timePeriodMap[tpID]; !exists {
			timePeriodMap[tpID] = &tenantmodels.ScheduleGroup{
				TimePeriod: tenantmodels.TimePeriod{
					ID:        tpID,
					SectionID: tpSectionID,
					StartTime: tpStartTime,
					EndTime:   tpEndTime,
				},
				Schedules: []tenantmodels.Schedule{},
			}
			timePeriodOrder = append(timePeriodOrder, tpID)
		}

		// Add schedule if it exists (not NULL)
		if sID.Valid {
			schedule := tenantmodels.Schedule{
				ID:            int(sID.Int64),
				SectionID:     int(sSectionID.Int64),
				TimePeriodID:  int(sTimePeriodID.Int64),
				SubjectCode:   sSubjectCode.String,
				SubjectName:   sSubjectName.String,
				Weekday:       sWeekday.String,
				ClassroomCode: sClassroomCode.String,
				Type:          sType.String,
				ColorConfig:   int(tColorConfig.Int64),
				TenantName:    tenantName.String,
				SectionName:   sectionName.String,
			}
			timePeriodMap[tpID].Schedules = append(timePeriodMap[tpID].Schedules, schedule)
		}
	}

	// Check for any iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Build the result in the correct order
	var result tenantmodels.ScheduleGroupCollection
	for _, tpID := range timePeriodOrder {
		result = append(result, *timePeriodMap[tpID])
	}

	return result, nil
}

// GetScheduleForSection TODO: Add description
func GetScheduleForSection(
	sectionID string,
	tenantDB *sql.DB,
) (tenantmodels.ScheduleGroupCollection, error) {
	query := `SELECT tp.id, tp.section_id, tp.start_time, tp.end_time,
    s.id, s.section_id, s.time_period_id, s.subject_code, sub.subject_name,
    s.weekday, s.classroom_code, s.type, NULL as color_config, NULL as tenant_name,
	NULL as section_name
    FROM time_periods tp
    LEFT JOIN schedule s ON tp.id = s.time_period_id AND tp.section_id = s.section_id
	LEFT JOIN ednevnik_workspace.subjects sub ON s.subject_code = sub.subject_code
    WHERE tp.section_id = ?
    ORDER BY tp.start_time`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return tenantmodels.ScheduleGroupCollection{}, err
	}

	return processScheduleRows(rows)
}

// GetScheduleForTeacher TODO: Add description
func GetScheduleForTeacher(
	teacherID string,
	tenantDB *sql.DB,
) (tenantmodels.ScheduleGroupCollection, error) {
	query := `SELECT tp.id, tp.section_id, tp.start_time, tp.end_time,
    s.id, s.section_id, s.time_period_id, s.subject_code, sub.subject_name,
    s.weekday, s.classroom_code, s.type, ten.color_config, ten.tenant_name,
	CONCAT('Odjeljenje ', sec.class_code, '-', sec.section_code) AS section_name
    FROM time_periods tp
    JOIN schedule s ON tp.id = s.time_period_id AND tp.section_id = s.section_id
    JOIN teachers_sections_subjects tss ON tp.section_id = tss.section_id
	JOIN sections sec ON tss.section_id = sec.id
	JOIN ednevnik_workspace.tenant ten ON sec.tenant_id = ten.id
	JOIN ednevnik_workspace.subjects sub ON s.subject_code = sub.subject_code
    AND s.subject_code = tss.subject_code
    WHERE tss.teacher_id = ? AND sec.archived = 0
    ORDER BY tp.start_time`

	rows, err := tenantDB.Query(query, teacherID)
	if err != nil {
		return tenantmodels.ScheduleGroupCollection{}, err
	}

	return processScheduleRows(rows)
}

// GetAllSchedulesForSection returns all schedules for section including historical ones, grouped by batch_id
func GetAllSchedulesForSection(
	sectionID int,
	tenantDB *sql.DB,
) ([]tenantmodels.ScheduleGroupCollection, error) {
	// Query to get all historical versions of time periods and schedules, grouped by batch_id
	query := `
	SELECT DISTINCT
		tp.batch_id, tp.row_start, tp.row_end, tp.id, tp.section_id, tp.start_time,
		tp.end_time, s.id, s.section_id, s.time_period_id, s.subject_code,
		sub.subject_name, s.weekday, s.classroom_code, s.type, s.row_start as schedule_row_start,
		s.row_end as schedule_row_end
	FROM time_periods FOR SYSTEM_TIME ALL tp
	JOIN schedule FOR SYSTEM_TIME ALL s ON tp.id = s.time_period_id
		AND tp.section_id = s.section_id
		AND tp.batch_id = s.batch_id  -- Ensure same batch
	LEFT JOIN ednevnik_workspace.subjects sub ON s.subject_code = sub.subject_code
	WHERE tp.section_id = ?
	ORDER BY tp.row_start ASC, tp.start_time ASC`

	rows, err := tenantDB.Query(query, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to group schedules by their batch_id
	schedulesByBatch := make(map[string]tenantmodels.ScheduleGroupCollection)
	batchOrder := []string{} // To maintain order

	for rows.Next() {
		var tp tenantmodels.TimePeriod
		var batchID string
		var tpRowStart, tpRowEnd, scheduleRowStart, scheduleRowEnd sql.NullString
		var scheduleID sql.NullInt64
		var timePeriodID, scheduleSectionID sql.NullInt64
		var subjectCode, subjectName, weekday, classroomCode, scheduleType sql.NullString

		err := rows.Scan(
			&batchID,
			&tpRowStart, &tpRowEnd,
			&tp.ID, &tp.SectionID, &tp.StartTime, &tp.EndTime,
			&scheduleID, &scheduleSectionID, &timePeriodID,
			&subjectCode, &subjectName, &weekday, &classroomCode, &scheduleType,
			&scheduleRowStart, &scheduleRowEnd,
		)
		if err != nil {
			return nil, err
		}

		// Convert weekday to Bosnian if needed
		if weekday.Valid && weekday.String != "" {
			if bosnianWeekday, exists := weekdayConvertToBosnianMap[weekday.String]; exists {
				weekday.String = bosnianWeekday
			}
		}

		// Initialize the schedule group collection for this batch if it doesn't exist
		if _, exists := schedulesByBatch[batchID]; !exists {
			schedulesByBatch[batchID] = tenantmodels.ScheduleGroupCollection{}
			batchOrder = append(batchOrder, batchID)
		}

		collection := schedulesByBatch[batchID]

		// Find the schedule group for this time period within this batch
		var scheduleGroup *tenantmodels.ScheduleGroup
		for i := range collection {
			if collection[i].TimePeriod.ID == tp.ID {
				scheduleGroup = &collection[i]
				break
			}
		}

		// Create schedule (always exists with JOIN)
		schedule := tenantmodels.Schedule{
			ID:            int(scheduleID.Int64),
			SectionID:     int(scheduleSectionID.Int64),
			TimePeriodID:  int(timePeriodID.Int64),
			SubjectCode:   subjectCode.String,
			SubjectName:   subjectName.String,
			Weekday:       weekday.String,
			ClassroomCode: classroomCode.String,
			Type:          scheduleType.String,
			RowStart:      scheduleRowStart.String,
			RowEnd:        scheduleRowEnd.String,
		}

		if scheduleGroup == nil {
			// Create new schedule group with this first schedule
			newGroup := tenantmodels.ScheduleGroup{
				TimePeriod: tp,
				Schedules:  []tenantmodels.Schedule{schedule},
				CreatedAt:  schedule.RowStart,
			}
			collection = append(collection, newGroup)
		} else {
			// Check if this exact schedule already exists in this group
			exists := false
			for _, existingSchedule := range scheduleGroup.Schedules {
				if existingSchedule.ID == schedule.ID {
					exists = true
					break
				}
			}

			if !exists {
				scheduleGroup.Schedules = append(scheduleGroup.Schedules, schedule)
			}
		}

		schedulesByBatch[batchID] = collection
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := make([]tenantmodels.ScheduleGroupCollection, 0, len(schedulesByBatch))
	for _, batchID := range batchOrder {
		if collection, exists := schedulesByBatch[batchID]; exists {
			result = append(result, collection)
		}
	}

	return result, nil
}
