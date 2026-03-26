package api

import (
	"fmt"
	"strconv"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
)

func resolveStudentID(db database.Database, schoolID *uint, value string) (uint, error) {
	if n, err := strconv.Atoi(value); err == nil {
		return uint(n), nil
	}
	var student models.Student
	if err := whereByIDOrUUID(db, value, schoolID).First(&student).Error(); err != nil {
		return 0, fmt.Errorf("student not found")
	}
	return student.ID, nil
}

func resolveClassID(db database.Database, schoolID *uint, value string) (uint, error) {
	if n, err := strconv.Atoi(value); err == nil {
		return uint(n), nil
	}
	var class models.Class
	if err := whereByIDOrUUID(db, value, schoolID).First(&class).Error(); err != nil {
		return 0, fmt.Errorf("class not found")
	}
	return class.ID, nil
}

func resolveEnrollmentID(db database.Database, schoolID *uint, value string) (uint, error) {
	if n, err := strconv.Atoi(value); err == nil {
		return uint(n), nil
	}
	var enrollment models.StudentEnrollment
	if err := whereByIDOrUUID(db, value, schoolID).First(&enrollment).Error(); err != nil {
		return 0, fmt.Errorf("enrollment not found")
	}
	return enrollment.ID, nil
}
