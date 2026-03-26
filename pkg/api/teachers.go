package api

import (
	"net/http"

	"golang-rest-api-template/pkg/models"

	"github.com/gin-gonic/gin"
)

func (r *userRepository) ListTeachers(c *gin.Context) {
	schoolID, _, ok := requireTenant(c)
	if !ok {
		return
	}
	offset, limit, ok := parsePagination(c)
	if !ok {
		return
	}

	query := r.DB.Model(&models.User{}).Where("role = ?", "guru")
	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if uuid := c.Query("uuid"); uuid != "" {
		query = query.Where("uuid = ?", uuid)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count teachers"})
		return
	}

	var teachers []models.User
	if err := query.Offset(offset).Limit(limit).Order("id asc").Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch teachers"})
		return
	}

	for i := range teachers {
		teachers[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": teachers, "meta": gin.H{"offset": offset, "limit": limit, "total": total, "count": len(teachers)}})
}
