package api

import (
	"fmt"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
)

func resolveEnrollmentForTerm(db database.Database, schoolID *uint, enrollmentID *uint, studentID uint, academicYear string, semester int) (*models.StudentEnrollment, error) {
	if schoolID == nil {
		return nil, fmt.Errorf("missing school context")
	}

	var enrollment models.StudentEnrollment
	if enrollmentID != nil {
		if err := db.Where("id = ? AND school_id = ?", *enrollmentID, *schoolID).First(&enrollment).Error(); err != nil {
			return nil, fmt.Errorf("enrollment not found")
		}
		if enrollment.StudentID != studentID {
			return nil, fmt.Errorf("enrollment does not belong to student")
		}
		if enrollment.AcademicYear != academicYear || enrollment.Semester != semester {
			return nil, fmt.Errorf("enrollment term mismatch")
		}
		return &enrollment, nil
	}

	if err := db.Where("student_id = ? AND school_id = ? AND academic_year = ? AND semester = ?", studentID, *schoolID, academicYear, semester).Order("id desc").First(&enrollment).Error; err != nil {
		return nil, fmt.Errorf("enrollment not found for student term")
	}
	return &enrollment, nil
}
