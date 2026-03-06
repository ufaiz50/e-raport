package api

import (
	"context"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type classRepository struct {
	DB  database.Database
	Ctx *context.Context
}

func NewClassRepository(db database.Database, ctx *context.Context) *classRepository {
	return &classRepository{DB: db, Ctx: ctx}
}

// FindClasses godoc
// @Summary Get all classes with pagination
// @Description Get a list of all classes with optional pagination
// @Tags classes
// @Security ApiKeyAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Class "Successfully retrieved list of classes"
// @Router /classes [get]
func (r *classRepository) FindClasses(c *gin.Context) {
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	var classes []models.Class
	r.DB.Offset(offset).Limit(limit).Order("name asc").Find(&classes)
	c.JSON(http.StatusOK, gin.H{"data": classes})
}

func (r *classRepository) CreateClass(c *gin.Context) {
	var input models.CreateClass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class := models.Class{Name: input.Name, Level: input.Level, Homeroom: input.Homeroom, AcademicYear: input.AcademicYear}
	if err := r.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create class"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": class})
}

func (r *classRepository) UpdateClass(c *gin.Context) {
	var class models.Class
	if err := r.DB.Where("id = ?", c.Param("id")).First(&class).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	var input models.UpdateClass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := r.DB.Model(&class).Updates(models.Class{Name: input.Name, Level: input.Level, Homeroom: input.Homeroom, AcademicYear: input.AcademicYear}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update class"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": class})
}

func (r *classRepository) DeleteClass(c *gin.Context) {
	var class models.Class
	if err := r.DB.Where("id = ?", c.Param("id")).First(&class).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	r.DB.Delete(&class)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
