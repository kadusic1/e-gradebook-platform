package tenantshared

import (
	"database/sql"
	commonmodels "ednevnik-backend/models/common"
	tenantmodels "ednevnik-backend/models/tenant"
	wpmodels "ednevnik-backend/models/workspace"
)

// ITenant TODO: Add description
type ITenant interface {
	// Function that creates tables in the database
	CreateSchema(tenantAdmin wpmodels.Teacher) error
	// Function that creates the database for the tenant
	CreateDB() (string, error)
	// Function that returns the prefix for the tenant database
	// e.g. "ednevnik_primary_" or "ednevnik_secondary_"
	GetDBPrefix() string
	// Function that deletes the tenant database
	DropDB() error
	// Function that return curriculum for the tenant
	GetCurriculumsForAssignment() ([]wpmodels.Curriculum, error)
	// Function that assigns curriculums to the tenant
	AssignCurriculumsToTenant(curriculumCodes []string) error
	GetCurriculumsForTenant() ([]wpmodels.CurriculumGet, error)
	UnassignCurriculumFromTenant(curriculumCode string) error
	GetDBName() (string, error)
	StartTransactions(workspaceDb *sql.DB, tenantDb *sql.DB) (*sql.Tx, *sql.Tx, error)
	GetTeachersForTenant() ([]wpmodels.Teacher, error)
	GetMetadataForSectionCreation() (*tenantmodels.SectionCreateMetadata, error)
	GrantTenantDBPrivileges() error
	RevokeTenantDBPrivileges() error
	GetSectionsForTenant(archived int) ([]tenantmodels.Section, error)
	DeleteTenantSection(tenantID string) error
	UpdateTenantSection(newSection tenantmodels.Section, sectionID string) (tenantmodels.Section, error)
	CreateTenantSection(newSection tenantmodels.SectionCreate) (tenantmodels.Section, error)
	GetPupilsForSection(sectionID string, includeUnenrolled bool) (*commonmodels.GetSectionPupilsResponse, error)
	DeletePupilFromSection(pupilID, sectionID string) error
	UpdatePupil(oldPupil tenantmodels.Pupil, newPupil tenantmodels.Pupil) error
	UpdateTenantSemesterDates(semesterCode, startDate, endDate, nppCode string) (wpmodels.TenantSemester, error)
	GetSemestersForTenant() ([]wpmodels.TenantSemester, error)
	SendPupilSectionInvite(pupilID int, sectionID, tenantID string) (*tenantmodels.PupilSectionInvite, error)
	GetPupilSectionInvite(inviteID int) (*tenantmodels.PupilSectionInvite, error)
	AcceptPupilSectionInvite(inviteID string) error
	DeclinePupilSectionInvite(inviteID string) error
	GetDataForTeacherInviteForTenant() ([]commonmodels.DataForTeacherSectionInvite, error)
	HandleTeacherSectionAssignments(
		teacherID int, assignmentRequest wpmodels.TeacherSectionAssignment,
	) ([]commonmodels.DataForTeacherSectionInvite, []wpmodels.Teacher, error)
	GetInvitesForTeacher(inviteID int) ([]wpmodels.TeacherSectionInviteRecord, error)
	GetAllTeacherInvites() ([]wpmodels.TeacherSectionInviteRecord, error)
	AcceptTeacherSectionInvite(inviteID string) error
	DeclineTeacherSectionInvite(inviteID string) error
	DeletePupilInvite(inviteID, pupilID int) error
	DeleteTeacherInvite(inviteID, teacherID int) error
	DeletePupilFromTenant(pupilID string) error
	DeleteTenantTeacherData(teacherID string) error
	DeleteTeacherFromTenant(teacherID string) error
	GetSectionsForTeacher(teacherID string, archived int) ([]tenantmodels.Section, error)
	GetSectionsForPupil(pupilID string, archived int) ([]tenantmodels.Section, error)
	CreateSchedule(data tenantmodels.ScheduleGroupCollection, sectionID string) error
	GetScheduleForSection(sectionID string) (tenantmodels.ScheduleGroupCollection, error)
	GetScheduleForTeacher(teacherID string) (tenantmodels.ScheduleGroupCollection, error)
	CreateClassroom(data tenantmodels.Classroom) error
	UpdateClassroom(data tenantmodels.Classroom, oldCode string) error
	GetAllClassroomsForTenant() ([]tenantmodels.Classroom, error)
	DeleteClassroom(code string) error
	GetLessonsForSection(sectionID int, claims *wpmodels.Claims) ([]tenantmodels.LessonData, error)
	CreateSectionLesson(requestData tenantmodels.LessonData, teacherID int) (*tenantmodels.LessonData, error)
	UpdateLesson(lessonID int, requestData tenantmodels.LessonData, teacherID int) (*tenantmodels.LessonData, error)
	DeleteLesson(lessonID int) error
	GetSubjectsForSection(sectionID int, claims *wpmodels.Claims) ([]wpmodels.Subject, error)
	GetLessonByID(lessonID int) (*tenantmodels.LessonData, error)
	GetAbsentAttendancesForSection(sectionID int) ([]tenantmodels.PupilAttendance, error)
	GetAbsentAttendancesForPupil(pupilID, sectionID int) ([]tenantmodels.PupilAttendance, error)
	HandleAttendanceAction(action tenantmodels.AttendanceAction) error
	GetPupilCountForSection(sectionID int) (int, error)
	GetSectionGradesForSubject(sectionID int, semesterCode, subjectCode string) ([]tenantmodels.GradePupilGroup, error)
	CreateGrade(grade *tenantmodels.Grade) (*tenantmodels.GradePupilGroup, error)
	DeleteGrade(grade *tenantmodels.Grade, teacherID int) (*tenantmodels.GradePupilGroup, error)
	UpdateGrade(grade *tenantmodels.Grade) (*tenantmodels.GradePupilGroup, error)
	GetPupilGradesForSectionPupil(sectionID, pupilID int, semesterCode string) ([]commonmodels.GradeSubjectGroup, error)
	UpdatePupilBehaviourGrade(behaviourGradesToUpdate tenantmodels.BehaviourGrade, teacherID int) (*tenantmodels.BehaviourGrade, error)
	GetSemestersForSection(sectionID string) ([]wpmodels.TenantSemester, error)
	GetSectionBehaviourGradesForPupil(pupilID, sectionID int) ([]tenantmodels.BehaviourGrade, error)
	ArchiveSection(sectionID int) error
	GetCertificateData(sectionID, pupilID int) (*commonmodels.Certificate, error)
	GetGradeEditHistory(gradeID int) ([]tenantmodels.Grade, error)
	GetBehaviourGradeHistory(behaviourGradeID int) ([]tenantmodels.BehaviourGrade, error)
	GetCompleteGradebookData(sectionID int) (*tenantmodels.CompleteGradebook, error)
	UnenrollPupilFromSection(pupilID, sectionID int) error
}
