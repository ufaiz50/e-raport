package api

import (
	"fmt"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
)

func resolveEnrollmentForTerm(db database.Database, schoolID *string, enrollmentID *string, studentID string, academicYear string, semester int) (*models.StudentEnrollment, error) {
	var enrollment models.StudentEnrollment

	base := db.Where("school_id = ?", *schoolID)
	if enrollmentID != nil && *enrollmentID != "" {
		if err := base.Where("id = ?", *enrollmentID).First(&enrollment).Error(); err == nil {
			return &enrollment, nil
		}
	}

	if err := base.Where("student_id = ? AND academic_year = ? AND semester = ?", studentID, academicYear, semester).First(&enrollment).Error(); err == nil {
		return &enrollment, nil
	}

	return nil, fmt.Errorf("enrollment for selected term not found")
}
