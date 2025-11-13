package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) HealthCheck(c *gin.Context) {
	sqlDB, err := h.pingDB.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "db_error"})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db_unreachable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
