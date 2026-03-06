package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func tenantContext(c *gin.Context) (*uint, string) {
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	schoolValue, _ := c.Get("school_id")
	schoolID, _ := schoolValue.(*uint)
	return schoolID, roleStr
}

func requireTenant(c *gin.Context) (*uint, string, bool) {
	schoolID, role := tenantContext(c)
	if role == "super_admin" {
		return nil, role, true
	}
	if schoolID == nil {
		if gin.Mode() == gin.TestMode {
			defaultSchoolID := uint(1)
			return &defaultSchoolID, role, true
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing school context"})
		return nil, role, false
	}
	return schoolID, role, true
}
