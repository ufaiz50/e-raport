package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang-rest-api-template/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type Database interface {
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
	Find(interface{}, ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) Database
	Delete(interface{}, ...interface{}) *gorm.DB
	Model(model interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) Database
	Updates(interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
	Error() error
}

type GormDatabase struct {
	*gorm.DB
}

func (db *GormDatabase) Where(query interface{}, args ...interface{}) Database {
	return &GormDatabase{db.DB.Where(query, args...)}
}

func (db *GormDatabase) First(dest interface{}, conds ...interface{}) Database {
	return &GormDatabase{db.DB.First(dest, conds...)}
}

func (db *GormDatabase) Error() error {
	return db.DB.Error
}

func NewDatabase() *gorm.DB {
	var database *gorm.DB
	var err error

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Error loading .env file")
	}

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")

	dbURl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_hostname, db_port, db_name)

	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(dbURl), &gorm.Config{})
		if err == nil {
			break
		} else {
			log.Printf("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}
	database.AutoMigrate(&models.School{})
	database.AutoMigrate(&models.AcademicYear{})
	database.AutoMigrate(&models.Semester{})
	database.AutoMigrate(&models.Curriculum{})
	database.AutoMigrate(&models.Subject{})
	database.AutoMigrate(&models.CurriculumSubject{})
	database.AutoMigrate(&models.Student{})
	database.AutoMigrate(&models.Class{})
	database.AutoMigrate(&models.StudentEnrollment{})
	database.AutoMigrate(&models.Teaching{})
	database.AutoMigrate(&models.Book{})
	database.AutoMigrate(&models.Grade{})
	database.AutoMigrate(&models.ReportCard{})
	database.AutoMigrate(&models.Attendance{})
	database.AutoMigrate(&models.ReportNote{})
	database.AutoMigrate(&models.User{})

	// Backfill schools columns from legacy school_profiles table when available.
	// Safe no-op if school_profiles table does not exist.
	database.Exec(`
		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1 FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = 'school_profiles'
			) THEN
				UPDATE schools s
				SET
					name = COALESCE(NULLIF(sp.school_name, ''), s.name),
					address = COALESCE(NULLIF(sp.address, ''), s.address),
					npsn = COALESCE(NULLIF(sp.npsn, ''), s.npsn),
					principal_name = COALESCE(NULLIF(sp.principal_name, ''), s.principal_name),
					principal_nip = COALESCE(NULLIF(sp.principal_nip, ''), s.principal_nip),
					headmaster_sign = COALESCE(NULLIF(sp.headmaster_sign, ''), s.headmaster_sign),
					school_stamp = COALESCE(NULLIF(sp.school_stamp, ''), s.school_stamp)
				FROM school_profiles sp
				WHERE sp.school_id = s.id;
			END IF;
		END $$;
	`)

	// Backfill initial active enrollment from legacy students.class_id when missing.
	database.Exec(`
		INSERT INTO student_enrollments (school_id, student_id, class_id, academic_year, semester, is_active, start_date, created_at, updated_at)
		SELECT s.school_id, s.id, s.class_id, COALESCE(c.academic_year, to_char(current_date, 'YYYY') || '/' || to_char(current_date + interval '1 year', 'YYYY')),
			1, true, NOW(), NOW(), NOW()
		FROM students s
		JOIN classes c ON c.id = s.class_id
		WHERE s.class_id IS NOT NULL
		AND NOT EXISTS (
			SELECT 1 FROM student_enrollments se
			WHERE se.student_id = s.id AND se.is_active = true
		);
	`)

	return database
}
