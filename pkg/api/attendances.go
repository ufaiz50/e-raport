package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type attendanceRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewAttendanceRepository(db database.Database, ctx *context.Context) *attendanceRepository {
	return &attendanceRepository{DB: db, Ctx: ctx}
}

func (r *attendanceRepository) FindAttendances(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	var attendances []models.Attendance
	query := r.DB.Model(&models.Attendance{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if academicYear := c.Query("academic_year"); academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count attendances"})
		return
	}

	query.Offset(offset).Limit(limit).Order("created_at desc").Find(&attendances)
	c.JSON(http.StatusOK, gin.H{
		"data": attendances,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(attendances),
		},
	})
}

func (r *attendanceRepository) UpsertAttendance(c *gin.Context) {
	var input models.UpsertAttendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", input.StudentID, *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var attendance models.Attendance
	enrollment, err := resolveEnrollmentForTerm(r.DB, schoolID, input.EnrollmentID, input.StudentID, input.AcademicYear, input.Semester)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.DB.Where("school_id = ? AND student_id = ? AND semester = ? AND academic_year = ?", *schoolID, input.StudentID, input.Semester, input.AcademicYear).First(&attendance).Error(); err != nil {
		attendance = models.Attendance{
			SchoolID:       schoolID,
			EnrollmentID:   &enrollment.ID,
			StudentID:      input.StudentID,
			Semester:       input.Semester,
			AcademicYear:   input.AcademicYear,
			SickDays:       input.SickDays,
			PermissionDays: input.PermissionDays,
			AbsentDays:     input.AbsentDays,
		}
		r.DB.Create(&attendance)
		c.JSON(http.StatusCreated, gin.H{"data": attendance})
		return
	}

	r.DB.Model(&attendance).Updates(models.Attendance{
		SchoolID:       schoolID,
		EnrollmentID:   &enrollment.ID,
		SickDays:       input.SickDays,
		PermissionDays: input.PermissionDays,
		AbsentDays:     input.AbsentDays,
	})
	c.JSON(http.StatusOK, gin.H{"data": attendance})
}
