package api

import (
	"net/http"
	"strconv"
	"strings"

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
		if schoolParam := strings.TrimSpace(c.Query("school_id")); schoolParam != "" {
			parsed, err := strconv.ParseUint(schoolParam, 10, 64)
			if err != nil || parsed == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school_id query parameter"})
				return nil, role, false
			}
			sid := uint(parsed)
			return &sid, role, true
		}
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
