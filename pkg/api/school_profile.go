package api

import (
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type schoolProfileRepository struct {
	DB database.Database
}

func NewSchoolProfileRepository(db database.Database) *schoolProfileRepository {
	return &schoolProfileRepository{DB: db}
}

func (r *schoolProfileRepository) Get(c *gin.Context) {
	if c.FullPath() == "/api/v1/school-profile" {
		c.Header("Deprecation", "true")
		c.Header("Sunset", "2026-12-31")
		c.Header("Link", "</api/v1/schools/profile>; rel=\"successor-version\"")
	}

	schoolID, role, ok := requireTenant(c)
	if !ok {
		return
	}

	var school models.School
	if role != "super_admin" {
		if err := r.DB.Where("id = ?", *schoolID).First(&school).Error(); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "school profile not found"})
			return
		}
	} else if schoolID != nil {
		if err := r.DB.Where("id = ?", *schoolID).First(&school).Error(); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "school profile not found"})
			return
		}
	} else if err := r.DB.Order("created_at asc").First(&school).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": models.UpsertSchoolProfile{
		SchoolID:       &school.ID,
		SchoolName:     school.Name,
		NPSN:           school.NPSN,
		Address:        school.Address,
		PrincipalName:  school.PrincipalName,
		PrincipalNIP:   school.PrincipalNIP,
		HeadmasterSign: school.HeadmasterSign,
		SchoolStamp:    school.SchoolStamp,
	}})
}

func (r *schoolProfileRepository) Upsert(c *gin.Context) {
	if c.FullPath() == "/api/v1/school-profile" {
		c.Header("Deprecation", "true")
		c.Header("Sunset", "2026-12-31")
		c.Header("Link", "</api/v1/schools/profile>; rel=\"successor-version\"")
	}

	var input models.UpsertSchoolProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolID, _, ok := resolveWriteSchoolID(c, input.SchoolID)
	if !ok {
		return
	}

	var school models.School
	if err := r.DB.Where("id = ?", *schoolID).First(&school).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school not found"})
		return
	}

	if err := r.DB.Model(&school).Updates(models.School{
		Name:           input.SchoolName,
		NPSN:           input.NPSN,
		Address:        input.Address,
		PrincipalName:  input.PrincipalName,
		PrincipalNIP:   input.PrincipalNIP,
		HeadmasterSign: input.HeadmasterSign,
		SchoolStamp:    input.SchoolStamp,
	}).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": school})
}
