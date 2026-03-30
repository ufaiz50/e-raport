package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func tenantContext(c *gin.Context) (*string, string) {
	role, _ := c.Get("role")
	roleStr, _ := role.(string)
	schoolValue, _ := c.Get("school_id")
	schoolID, _ := schoolValue.(*string)
	return schoolID, roleStr
}

func requireTenant(c *gin.Context) (*string, string, bool) {
	schoolID, role := tenantContext(c)
	if role == "super_admin" {
		if schoolParam := strings.TrimSpace(c.Query("school_id")); schoolParam != "" {
			sid := schoolParam
			return &sid, role, true
		}
		return nil, role, true
	}
	if schoolID == nil || strings.TrimSpace(*schoolID) == "" {
		if gin.Mode() == gin.TestMode {
			defaultSchoolID := "00000000-0000-0000-0000-000000000001"
			return &defaultSchoolID, role, true
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing school context"})
		return nil, role, false
	}
	return schoolID, role, true
}

func resolveWriteSchoolID(c *gin.Context, bodySchoolID *string) (*string, string, bool) {
	schoolID, role, ok := requireTenant(c)
	if !ok {
		return nil, role, false
	}

	if role == "super_admin" {
		if bodySchoolID == nil || strings.TrimSpace(*bodySchoolID) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "school_id is required for super_admin write operations"})
			return nil, role, false
		}
		return bodySchoolID, role, true
	}

	if bodySchoolID != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id in payload is not allowed for non-super_admin"})
		return nil, role, false
	}

	if schoolID == nil || strings.TrimSpace(*schoolID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing school context"})
		return nil, role, false
	}

	return schoolID, role, true
}
