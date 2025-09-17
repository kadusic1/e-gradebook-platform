package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var firstNames = []string{"Ajdin", "Sara", "Eldar", "Samra", "Muhamed", "Hana", "Lejla", "Tarik", "Adis", "Jasmina"}
var lastNames = []string{"Ahmetovic", "Delic", "Begovic", "Hodzic", "Spahic", "Ibric", "Kavazovic", "Mujic", "Kadric", "Softic"}

type TenantData struct {
	ID   string
	Name string
}

type CurriculumData struct {
	CurriculumCode string
	ClassCode      string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var err error

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/ednevnik_workspace",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("MARIADB_HOST"),
		os.Getenv("MARIADB_PORT"),
	)
	var workspaceDB *sql.DB
	workspaceDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer workspaceDB.Close()

	tenantCommand := `SELECT id, tenant_name FROM tenant`
	rows, err := workspaceDB.Query(tenantCommand)
	if err != nil {
		log.Fatalf("Error querying tenant data: %v", err)
	}
	defer rows.Close()

	var tenants []TenantData
	for rows.Next() {
		var tenant TenantData
		if err := rows.Scan(&tenant.ID, &tenant.Name); err != nil {
			log.Fatalf("Error scanning tenant data: %v", err)
		}
		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error with rows: %v", err)
	}

	fmt.Println("Choose a tenant by ID:")
	fmt.Println(strings.Repeat("-", 40))
	for _, tenant := range tenants {
		fmt.Printf("ID: %s Name: %s\n", tenant.ID, tenant.Name)
	}
	fmt.Println(strings.Repeat("-", 40))
	var tenantID string
	fmt.Print("Enter Tenant ID: ")
	fmt.Scanln(&tenantID)

	// Check if tenant exists
	var tenantExists bool
	for _, tenant := range tenants {
		if tenant.ID == tenantID {
			tenantExists = true
			break
		}
	}

	if !tenantExists {
		log.Fatalf("Tenant with ID %s does not exist.", tenantID)
	}

	var tenantCurriculums []CurriculumData

	curriculumQuery := `SELECT ct.curriculum_code, c.class_code FROM
	curriculum_tenant ct
	JOIN curriculum c ON ct.curriculum_code = c.curriculum_code
	WHERE tenant_id = ?`

	rows, err = workspaceDB.Query(curriculumQuery, tenantID)
	if err != nil {
		log.Fatalf("Error querying curriculum data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var curriculum CurriculumData
		if err := rows.Scan(&curriculum.CurriculumCode, &curriculum.ClassCode); err != nil {
			log.Fatalf("Error scanning curriculum data: %v", err)
		}
		tenantCurriculums = append(tenantCurriculums, curriculum)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error with rows: %v", err)
	}

	if len(tenantCurriculums) == 0 {
		log.Fatalf("No curriculums found for tenant ID %s.", tenantID)
	}

	// Make a tenant DB connection
	tenantDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/ednevnik_tenant_db_tenant_id_%s",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASSWORD"),
		os.Getenv("MARIADB_HOST"),
		os.Getenv("MARIADB_PORT"),
		tenantID,
	)
	tenantDB, err := sql.Open("mysql", tenantDSN)
	if err != nil {
		log.Fatalf("Error connecting to tenant database: %v", err)
	}
	defer tenantDB.Close()

	for _, curriculum := range tenantCurriculums {
		for _, sectionCode := range []string{"A", "B"} {
			year := "2025/2026"
			res, err := tenantDB.Exec(
				`INSERT INTO sections (section_code, class_code, year, tenant_id,
			curriculum_code) VALUES (?, ?, ?, ?, ?)`,
				sectionCode, curriculum.ClassCode, year, tenantID, curriculum.CurriculumCode,
			)
			if err != nil {
				log.Fatalf("Error inserting section: %v", err)
			}

			sectionID, err := res.LastInsertId()
			if err != nil {
				log.Fatalf("Error getting section ID: %v", err)
			}

			// Fetch subjects for curriculum
			subjectRows, err := workspaceDB.Query(
				`SELECT subject_code FROM curriculum_subjects WHERE curriculum_code = ?`,
				curriculum.CurriculumCode,
			)
			if err != nil {
				log.Fatalf("Error fetching subjects: %v", err)
			}
			var subjectCodes []string
			for subjectRows.Next() {
				var subjectCode string
				if err := subjectRows.Scan(&subjectCode); err != nil {
					log.Fatalf("Error scanning subject_code: %v", err)
				}
				subjectCodes = append(subjectCodes, subjectCode)
			}
			subjectRows.Close()

			// Fetch semesters for curriculum (example: 1POL, 2POL)
			semesterRows, err := workspaceDB.Query(
				`SELECT semester_code FROM npp_semester WHERE
				npp_code = (SELECT npp_code FROM curriculum WHERE curriculum_code = ?)`,
				curriculum.CurriculumCode,
			)
			if err != nil {
				log.Fatalf("Error fetching semesters: %v", err)
			}
			var semesterCodes []string
			for semesterRows.Next() {
				var semesterCode string
				if err := semesterRows.Scan(&semesterCode); err != nil {
					log.Fatalf("Error scanning semester_code: %v", err)
				}
				semesterCodes = append(semesterCodes, semesterCode)
			}
			semesterRows.Close()

			batchID := uuid.NewString()
			res, err = tenantDB.Exec(
				`INSERT INTO time_periods (section_id, start_time, end_time, batch_id)
            VALUES (?, ?, ?, ?)`,
				sectionID, "08:00", "08:45", batchID,
			)
			if err != nil {
				log.Fatalf("Error inserting time period: %v", err)
			}
			timePeriodID, err := res.LastInsertId()
			if err != nil {
				log.Fatalf("Error getting time period ID: %v", err)
			}

			// Insert a sample schedule for the section (use first subject)
			if len(subjectCodes) > 0 {
				_, err = tenantDB.Exec(
					`INSERT INTO schedule (section_id, time_period_id,
					subject_code, weekday, type, batch_id)
                	VALUES (?, ?, ?, ?, ?, ?)`,
					sectionID, timePeriodID, subjectCodes[0], "Monday",
					"regular", batchID,
				)
				if err != nil {
					log.Fatalf("Error inserting schedule: %v", err)
				}
			}

			// 2. Create 5 pupils and enroll them
			for i := 1; i <= 5; i++ {
				// Generate random names
				pupilName := firstNames[rand.Intn(len(firstNames))]
				pupilLastName := lastNames[rand.Intn(len(lastNames))]

				email := fmt.Sprintf("%s.%s%d@gmail.com", strings.ToLower(pupilName), strings.ToLower(pupilLastName), rand.Intn(10000))

				log.Printf(
					"Creating pupil: %d for section: %s-%s",
					i,
					curriculum.ClassCode,
					sectionCode,
				)
				// Create an account for the pupil
				res, err = workspaceDB.Exec(
					`INSERT INTO accounts (email, password, account_type) VALUES (?, ?, ?)`,
					email,
					"$2a$10$CKRfsoJMZcpHxXgsU8b/du7JdZd2IHgCEAqbOhgL0L3fayRWFDHyi",
					"pupil",
				)
				if err != nil {
					log.Fatalf("Error inserting account: %v", err)
				}
				accountID, err := res.LastInsertId()
				if err != nil {
					log.Fatalf("Error getting account ID: %v", err)
				}

				phoneNumber := fmt.Sprintf("06%s%s%s%02d", sectionCode, curriculum.ClassCode, tenantID, i)
				guardianNumber := fmt.Sprintf("06%s%s%s%02d", sectionCode, curriculum.ClassCode, tenantID, i+20)

				// Insert pupil into workspace pupil_global
				res, err := workspaceDB.Exec(
					`INSERT INTO pupil_global(name, last_name, gender, address,
				guardian_name, phone_number, guardian_number, date_of_birth,
				religion, account_id, place_of_birth
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					pupilName, pupilLastName, "M", "Ulica", "Staratelj",
					phoneNumber, guardianNumber, "2005-01-01", "Islam",
					accountID, "Grad",
				)
				if err != nil {
					log.Fatalf("Error inserting pupil_global: %v", err)
				}
				pupilID, err := res.LastInsertId()
				if err != nil {
					log.Fatalf("Error getting pupil_global ID: %v", err)
				}

				// Insert pupil into tenant pupils table
				_, err = tenantDB.Exec(
					`INSERT INTO pupils(id, name, last_name, gender, address,
				guardian_name, phone_number, guardian_number, date_of_birth,
				religion, account_id, place_of_birth
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
					pupilID, pupilName, pupilLastName, "M", "Ulica", "Staratelj",
					phoneNumber, guardianNumber, "2005-01-01", "Islam",
					accountID, "Grad",
				)
				if err != nil {
					log.Fatalf("Error inserting pupil: %v", err)
				}

				// Enroll pupil in section
				_, err = tenantDB.Exec(
					`INSERT INTO pupils_sections (pupil_id, section_id) VALUES (?, ?)`,
					pupilID, sectionID,
				)
				if err != nil {
					log.Fatalf("Error enrolling pupil in section: %v", err)
				}
				_, err = workspaceDB.Exec(
					`INSERT INTO pupil_tenant (pupil_id, tenant_id) VALUES (?, ?)`,
					pupilID, tenantID,
				)
				if err != nil {
					log.Fatalf("Error enrolling pupil in tenant: %v", err)
				}

				// Insert one lesson for the section
				res, err = tenantDB.Exec(
					`INSERT INTO class_lesson (description, date, period_number,
					section_id, subject_code, signature)
					VALUES (?, ?, ?, ?, ?, ?)`,
					fmt.Sprintf(
						"Demo čas za %s-%s", curriculum.ClassCode, sectionCode,
					),
					"2025-09-01",
					1,
					sectionID,
					func() string {
						if len(subjectCodes) > 0 {
							return subjectCodes[pupilID%int64(len(subjectCodes))]
						}
						return ""
					}(),
					"Script User",
				)
				if err != nil {
					log.Fatalf("Error inserting demo lesson: %v", err)
				}
				lessonID, err := res.LastInsertId()
				if err != nil {
					log.Fatalf("Error getting demo lesson ID: %v", err)
				}

				var absentStatus = []string{"absent", "unexcused", "excused"}

				// Insert one absence for the first pupil for the lesson
				_, err = tenantDB.Exec(
					`INSERT INTO pupil_attendance (pupil_id, lesson_id, status)
					VALUES (?, ?, ?)`,
					pupilID, lessonID, absentStatus[pupilID%3],
				)
				if err != nil {
					log.Fatalf("Error inserting demo absence: %v", err)
				}

				var values []string
				var args []interface{}

				// 3. Insert grades for all subjects in curriculum for each semester
				// Optimization: bulk insert
				for _, subjectCode := range subjectCodes {
					for _, semesterCode := range semesterCodes {
						// Insert exam and final grades
						for _, gradeType := range []string{"exam", "final"} {
							grade := 2 + (i % 4)
							values = append(values, "(?, ?, ?, ?, ?, ?, ?)")
							args = append(
								args, gradeType, pupilID, sectionID, subjectCode,
								grade, semesterCode, "Script User",
							)
						}
					}
				}

				gradeInsertSql := fmt.Sprintf(
					`INSERT INTO student_grades (type, pupil_id, section_id,
				subject_code, grade, semester_code, signature) VALUES %s`,
					strings.Join(values, ", "),
				)
				_, err = tenantDB.Exec(gradeInsertSql, args...)
				if err != nil {
					log.Fatalf("Error inserting grades: %v", err)
				}
			}
		}
	}
}
