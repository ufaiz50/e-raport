package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePagination(c *gin.Context) (int, int, bool) {
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset format"})
		return 0, 0, false
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit format"})
		return 0, 0, false
	}

	return offset, limit, true
}
