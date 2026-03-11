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
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.Class{})
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count classes"})
		return
	}

	var classes []models.Class
	dataQuery := r.DB.Offset(offset).Limit(limit).Order("name asc")
	if schoolID != nil {
		dataQuery = dataQuery.Where("school_id = ?", *schoolID)
	}
	dataQuery.Find(&classes)
	c.JSON(http.StatusOK, gin.H{
		"data": classes,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(classes),
		},
	})
}

func (r *classRepository) CreateClass(c *gin.Context) {
	var input models.CreateClass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}
	class := models.Class{Name: input.Name, Level: input.Level, Homeroom: input.Homeroom, AcademicYear: input.AcademicYear, SchoolID: schoolID}
	if err := r.DB.Create(&class).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create class"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": class})
}

func (r *classRepository) UpdateClass(c *gin.Context) {
	var input models.UpdateClass
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&class).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	if err := r.DB.Model(&class).Updates(models.Class{Name: input.Name, Level: input.Level, Homeroom: input.Homeroom, AcademicYear: input.AcademicYear, SchoolID: schoolID}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update class"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": class})
}

func (r *classRepository) DeleteClass(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}

	var class models.Class
	if err := r.DB.Where("id = ? AND school_id = ?", c.Param("id"), *schoolID).First(&class).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "class not found"})
		return
	}
	r.DB.Delete(&class)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
