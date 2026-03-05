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
	var attendances []models.Attendance
	query := r.DB
	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if academicYear := c.Query("academic_year"); academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	query.Order("id desc").Find(&attendances)
	c.JSON(http.StatusOK, gin.H{"data": attendances})
}

func (r *attendanceRepository) UpsertAttendance(c *gin.Context) {
	var input models.UpsertAttendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ?", input.StudentID).First(&student).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
		return
	}

	var attendance models.Attendance
	if err := r.DB.Where("student_id = ? AND semester = ? AND academic_year = ?", input.StudentID, input.Semester, input.AcademicYear).First(&attendance).Error(); err != nil {
		attendance = models.Attendance{
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
		SickDays:       input.SickDays,
		PermissionDays: input.PermissionDays,
		AbsentDays:     input.AbsentDays,
	})
	c.JSON(http.StatusOK, gin.H{"data": attendance})
}
