package api

import (
	"context"
	"net/http"
	"time"

	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

type enrollmentRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewEnrollmentRepository(db database.Database, ctx *context.Context) *enrollmentRepository {
	return &enrollmentRepository{DB: db, Ctx: ctx}
}

func (r *enrollmentRepository) FindEnrollments(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.StudentEnrollment{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if studentID := c.Query("student_id"); studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	if academicYear := c.Query("academic_year"); academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}
	if semester := c.Query("semester"); semester != "" {
		query = query.Where("semester = ?", semester)
	}
	if active := c.Query("is_active"); active != "" {
		query = query.Where("is_active = ?", active == "true")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count enrollments"})
		return
	}

	var rows []models.StudentEnrollment
	if err := query.Offset(offset).Limit(limit).Order("id desc").Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch enrollments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rows, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(rows)}})
}

func (r *enrollmentRepository) CreateEnrollment(c *gin.Context) {
	var input models.CreateEnrollment
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
	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", input.ClassID, *schoolID).First(&class).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
		return
	}

	now := time.Now()
	r.DB.Model(&models.StudentEnrollment{}).
		Where("student_id = ? AND school_id = ? AND is_active = ?", input.StudentID, *schoolID, true).
		Updates(map[string]interface{}{"is_active": false, "end_date": now})

	startDate := now
	if input.StartDate != nil {
		startDate = *input.StartDate
	}

	enrollment := models.StudentEnrollment{
		SchoolID:     schoolID,
		StudentID:    input.StudentID,
		ClassID:      input.ClassID,
		AcademicYear: input.AcademicYear,
		Semester:     input.Semester,
		IsActive:     true,
		StartDate:    startDate,
	}
	if err := r.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create enrollment"})
		return
	}

	r.DB.Model(&student).Updates(models.Student{ClassID: &input.ClassID})
	c.JSON(http.StatusCreated, gin.H{"data": enrollment})
}

func (r *enrollmentRepository) CloseEnrollment(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var enrollment models.StudentEnrollment
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&enrollment).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
		return
	}

	now := time.Now()
	if err := r.DB.Model(&enrollment).Updates(map[string]interface{}{"is_active": false, "end_date": now}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to close enrollment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": enrollment})
}
