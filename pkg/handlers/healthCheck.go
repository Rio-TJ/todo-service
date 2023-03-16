package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": "todo",
		"version": "0.0.0",
		"status":  "OK",
	})
}
