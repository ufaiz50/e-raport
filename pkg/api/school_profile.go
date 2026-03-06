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
	var profile models.SchoolProfile
	if err := r.DB.Order("id asc").First(&profile).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "school profile not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": profile})
}

func (r *schoolProfileRepository) Upsert(c *gin.Context) {
	var input models.UpsertSchoolProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.SchoolProfile
	if err := r.DB.Order("id asc").First(&profile).Error; err != nil {
		profile = models.SchoolProfile{
			SchoolID:       input.SchoolID,
			SchoolName:     input.SchoolName,
			NPSN:           input.NPSN,
			Address:        input.Address,
			PrincipalName:  input.PrincipalName,
			PrincipalNIP:   input.PrincipalNIP,
			HeadmasterSign: input.HeadmasterSign,
			SchoolStamp:    input.SchoolStamp,
		}
		if err := r.DB.Create(&profile).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create profile"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"data": profile})
		return
	}

	if err := r.DB.Model(&profile).Updates(models.SchoolProfile{
		SchoolID:       input.SchoolID,
		SchoolName:     input.SchoolName,
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

	c.JSON(http.StatusOK, gin.H{"data": profile})
}
