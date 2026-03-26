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
	query := r.DB.Model(&models.School{})
	if uuid := c.Query("uuid"); uuid != "" {
		query = query.Where("uuid = ?", uuid)
	}
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count schools"})
		return
	}

	var schools []models.School
	query.Offset(offset).Limit(limit).Order("id asc").Find(&schools)

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

	school := models.School{
		Name:           input.Name,
		Code:           input.Code,
		Address:        input.Address,
		NPSN:           input.NPSN,
		PrincipalName:  input.PrincipalName,
		PrincipalNIP:   input.PrincipalNIP,
		HeadmasterSign: input.HeadmasterSign,
		SchoolStamp:    input.SchoolStamp,
	}
	if err := r.DB.Create(&school).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create school"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": school})
}

// UpdateSchool godoc
// @Summary Update school
// @Description Update school by id (super admin only)
// @Tags schools
// @Security ApiKeyAuth
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param id path int true "School ID"
// @Param input body models.School true "School object"
// @Success 200 {object} models.School
// @Router /schools/{id} [put]
func (r *schoolRepository) UpdateSchool(c *gin.Context) {
	var school models.School
	if err := whereByIDOrUUID(r.DB, c.Param("id"), nil).First(&school).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	var input models.School
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.DB.Model(&school).Updates(models.School{
		Name:           input.Name,
		Code:           input.Code,
		Address:        input.Address,
		NPSN:           input.NPSN,
		PrincipalName:  input.PrincipalName,
		PrincipalNIP:   input.PrincipalNIP,
		HeadmasterSign: input.HeadmasterSign,
		SchoolStamp:    input.SchoolStamp,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update school"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": school})
}

// DeleteSchool godoc
// @Summary Delete school
// @Description Delete school by id (super admin only)
// @Tags schools
// @Security ApiKeyAuth
// @Security JwtAuth
// @Produce json
// @Param id path int true "School ID"
// @Success 204 {string} string "Successfully deleted school"
// @Router /schools/{id} [delete]
func (r *schoolRepository) DeleteSchool(c *gin.Context) {
	var school models.School
	if err := whereByIDOrUUID(r.DB, c.Param("id"), nil).First(&school).Error(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	// guard: prevent deleting school with dependent data
	var count int64
	r.DB.Where("school_id = ?", school.ID).Model(&models.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete school: users still linked"})
		return
	}
	r.DB.Where("school_id = ?", school.ID).Model(&models.Student{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete school: students still linked"})
		return
	}
	r.DB.Where("school_id = ?", school.ID).Model(&models.Class{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete school: classes still linked"})
		return
	}

	r.DB.Delete(&school)
	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
