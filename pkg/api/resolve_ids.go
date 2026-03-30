package api

import (
	"fmt"
	"strings"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
)

func resolveStudentID(db database.Database, schoolID *string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("student not found")
	}
	var student models.Student
	if err := whereByIDOrUUID(db, value, schoolID).First(&student).Error(); err != nil {
		return "", fmt.Errorf("student not found")
	}
	return student.ID, nil
}

func resolveClassID(db database.Database, schoolID *string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("class not found")
	}
	var class models.Class
	if err := whereByIDOrUUID(db, value, schoolID).First(&class).Error(); err != nil {
		return "", fmt.Errorf("class not found")
	}
	return class.ID, nil
}

func resolveEnrollmentID(db database.Database, schoolID *string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("enrollment not found")
	}
	var enrollment models.StudentEnrollment
	if err := whereByIDOrUUID(db, value, schoolID).First(&enrollment).Error(); err != nil {
		return "", fmt.Errorf("enrollment not found")
	}
	return enrollment.ID, nil
}
