package api

import (
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type schoolRepository struct {
	DB database.Database
}

func NewSchoolRepository(db database.Database) *schoolRepository {
	return &schoolRepository{DB: db}
}

// ListSchools godoc
// @Summary Get all schools
// @Description Get all schools (super admin only)
// @Tags schools
// @Security ApiKeyAuth
// @Security JwtAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.School
// @Router /schools [get]
func (r *schoolRepository) ListSchools(c *gin.Context) {
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	var total int64
	if err := r.DB.Model(&models.School{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count schools"})
		return
	}

	var schools []models.School
	r.DB.Offset(offset).Limit(limit).Order("id asc").Find(&schools)

	c.JSON(http.StatusOK, gin.H{
		"data": schools,
		"meta": gin.H{
			"offset": offset,
			"limit":  limit,
			"total":  total,
			"count":  len(schools),
		},
	})
}

// CreateSchool godoc
// @Summary Create school
// @Description Create a new school (super admin only)
// @Tags schools
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param input body models.School true "School object"
// @Success 201 {object} models.School
// @Router /schools [post]
func (r *schoolRepository) CreateSchool(c *gin.Context) {
	var input models.School
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school := models.School{Name: input.Name, Code: input.Code, Address: input.Address}
	if err := r.DB.Create(&school).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create school"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": school})
}
