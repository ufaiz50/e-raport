package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentRepository interface {
	FindStudents(c *gin.Context)
	CreateStudent(c *gin.Context)
	FindStudent(c *gin.Context)
	UpdateStudent(c *gin.Context)
	DeleteStudent(c *gin.Context)
}

type studentRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func (r *studentRepository) syncActiveEnrollment(student *models.Student, classID *uint, schoolID *uint) error {
	if classID == nil {
		return nil
	}

	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", *classID, *schoolID).First(&class).Error(); err != nil {
		return err
	}

	var active models.StudentEnrollment
	if err := r.DB.Where("student_id = ? AND school_id = ? AND is_active = ?", student.ID, *schoolID, true).Order("id desc").First(&active).Error; err == nil {
		if active.ClassID == *classID && active.AcademicYear == class.AcademicYear {
			return nil
		}
		now := time.Now()
		r.DB.Model(&active).Updates(map[string]interface{}{"is_active": false, "end_date": now})
	}

	enrollment := models.StudentEnrollment{
		SchoolID:     schoolID,
		StudentID:    student.ID,
		ClassID:      *classID,
		AcademicYear: class.AcademicYear,
		Semester:     1,
		IsActive:     true,
		StartDate:    time.Now(),
	}
	return r.DB.Create(&enrollment).Error
}

func NewStudentRepository(db database.Database, ctx *context.Context) *studentRepository {
	return &studentRepository{DB: db, Ctx: ctx}
}

// FindStudents godoc
// @Summary Get all students with pagination
// @Description Get a list of all students with optional pagination
// @Tags students
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Student "Successfully retrieved list of students"
// @Router /students [get]
func (r *studentRepository) FindStudents(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Student{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count students"})
		return
	}

	var students []models.Student
	dataQuery := r.DB.Offset(offset).Limit(limit).Order("id asc")
	if schoolID != nil {
		dataQuery = dataQuery.Where("school_id = ?", *schoolID)
	}
	dataQuery.Find(&students)
	for i := range students {
		if students[i].SchoolID == nil {
			continue
		}
		var active models.StudentEnrollment
		if err := r.DB.Where("student_id = ? AND school_id = ? AND is_active = ?", students[i].ID, *students[i].SchoolID, true).Order("id desc").First(&active).Error; err == nil {
			students[i].ClassID = &active.ClassID
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data": students,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(students),
		},
	})
}

// CreateStudent godoc
// @Summary Create a new student
// @Description Create a new student with kindergarten type (junior/senior)
// @Tags students
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param input body models.CreateStudent true "Create student object"
// @Success 201 {object} models.Student "Successfully created student"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /students [post]
func (r *studentRepository) CreateStudent(c *gin.Context) {
	var input models.CreateStudent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	student := models.Student{Name: input.Name, Email: input.Email, Type: input.Type, SchoolID: schoolID}
	if input.ClassID != nil {
		var class models.Class
		if err := r.DB.Where("id = ? AND school_id = ?", *input.ClassID, *schoolID).First(&class).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
			return
		}
		student.ClassID = input.ClassID
	}
	r.DB.Create(&student)
	if input.ClassID != nil {
		if err := r.syncActiveEnrollment(&student, input.ClassID, schoolID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create enrollment"})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"data": student})
}

// FindStudent godoc
// @Summary Find a student by ID
// @Description Get details of a student by ID
// @Tags students
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 200 {object} models.Student "Successfully retrieved student"
// @Failure 404 {string} string "Student not found"
// @Router /students/{id} [get]
func (r *studentRepository) FindStudent(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

// UpdateStudent godoc
// @Summary Update a student by ID
// @Description Update student details and kindergarten type (junior/senior)
// @Tags students
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID"
// @Param input body models.UpdateStudent true "Update student object"
// @Success 200 {object} models.Student "Successfully updated student"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Student not found"
// @Router /students/{id} [put]
func (r *studentRepository) UpdateStudent(c *gin.Context) {
	var input models.UpdateStudent

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	if input.ClassID != nil {
		var class models.Class
		if err := r.DB.Where("id = ? AND school_id = ?", *input.ClassID, *schoolID).First(&class).Error(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "class not found"})
			return
		}
	}

	r.DB.Model(&student).Updates(models.Student{Name: input.Name, Email: input.Email, Type: input.Type, SchoolID: schoolID, ClassID: input.ClassID})
	if input.ClassID != nil {
		if err := r.syncActiveEnrollment(&student, input.ClassID, schoolID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update enrollment"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

// DeleteStudent godoc
// @Summary Delete a student by ID
// @Description Delete student by ID
// @Tags students
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Student ID"
// @Success 204 {string} string "Successfully deleted student"
// @Failure 404 {string} string "Student not found"
// @Router /students/{id} [delete]
func (r *studentRepository) DeleteStudent(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var student models.Student
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	r.DB.Delete(&student)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
