package main

import (
	"database/sql"
	"ednevnik-backend/util"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	var err error

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := util.BuildDBConnectionString("ednevnik_workspace")

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Transaction start error: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Fatalf("Transaction rollback error: %v", err)
		}
	}()

	for _, canton := range Cantons {
		if err = util.InsertCanton(tx, canton); err != nil {
			log.Fatalf("Canton error: %v", err)
		}
	}

	for _, class := range Classes {
		if err = util.InsertClass(tx, class); err != nil {
			log.Fatalf("Class error: %v", err)
		}
	}

	for _, subject := range Subjects {
		if err = util.InsertSubject(tx, subject); err != nil {
			log.Fatalf("Subject error: %v", err)
		}
	}

	for _, npp := range Npps {
		if err = util.InsertNPP(tx, npp); err != nil {
			log.Fatalf("NPP error: %v", err)
		}
	}

	for _, course := range Courses {
		if err = util.InsertCourse(tx, course); err != nil {
			log.Fatalf("Course error: %v", err)
		}
	}

	for _, semester := range Semesters {
		if err = util.InsertSemester(tx, semester); err != nil {
			log.Fatalf("Semester error: %v", err)
		}
	}

	for _, nppSemester := range NPPSemester {
		if err = util.InsertNPPSemester(tx, nppSemester); err != nil {
			log.Fatalf("NPP Semester error: %v", err)
		}
	}

	for _, curriculum := range CurriculumsZDK {
		if err = util.InsertCurriculum(tx, curriculum); err != nil {
			log.Fatalf("Curriculum error: %v", err)
		}
	}

	for _, curriculum := range CurriculumsSBK {
		if err = util.InsertCurriculum(tx, curriculum); err != nil {
			log.Fatalf("Curriculum error: %v", err)
		}
	}

	for _, curriculumSubject := range CurriculumSubjectsZDKPrimary {
		if err = util.InsertCurriculumSubject(tx, curriculumSubject); err != nil {
			log.Fatalf("Curriculum subject error: %v", err)
		}
	}

	for _, curriculumSubject := range CurriculumSubjectsSBKPrimary {
		if err = util.InsertCurriculumSubject(tx, curriculumSubject); err != nil {
			log.Fatalf("Curriculum subject error: %v", err)
		}
	}

	for _, curriculumSubject := range CurriculumSubjectsZDKRTIASecondary {
		if err = util.InsertCurriculumSubject(tx, curriculumSubject); err != nil {
			log.Fatalf("Curriculum subject error: %v", err)
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Fatalf("Transaction commit error: %v", err)
	}

	log.Println("Static data inserted successfully!")
}
