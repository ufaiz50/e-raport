package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

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
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var input models.CreateStudent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var student models.Student
	var input models.UpdateStudent

	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
