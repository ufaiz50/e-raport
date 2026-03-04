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

func (r *studentRepository) FindStudents(c *gin.Context) {
	var students []models.Student
	r.DB.Find(&students)
	c.JSON(http.StatusOK, gin.H{"data": students})
}

func (r *studentRepository) CreateStudent(c *gin.Context) {
	var input models.CreateStudent
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student := models.Student{Name: input.Name, Email: input.Email}
	r.DB.Create(&student)
	c.JSON(http.StatusCreated, gin.H{"data": student})
}

func (r *studentRepository) FindStudent(c *gin.Context) {
	var student models.Student
	if err := r.DB.Where("id = ?", c.Param("id")).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (r *studentRepository) UpdateStudent(c *gin.Context) {
	var student models.Student
	var input models.UpdateStudent

	if err := r.DB.Where("id = ?", c.Param("id")).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r.DB.Model(&student).Updates(models.Student{Name: input.Name, Email: input.Email})
	c.JSON(http.StatusOK, gin.H{"data": student})
}

func (r *studentRepository) DeleteStudent(c *gin.Context) {
	var student models.Student
	if err := r.DB.Where("id = ?", c.Param("id")).First(&student).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	r.DB.Delete(&student)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
