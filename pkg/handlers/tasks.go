package handlers

import (
	"github.com/SoftclubIT/todo-service/pkg/models"
	"github.com/SoftclubIT/todo-service/pkg/scopes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	queryParams := c.Request.URL.Query()

	result := h.DB.Scopes(scopes.Paginate(queryParams)).Find(&tasks)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	paginationData, err := models.NewPaginationData(h.DB, &models.Task{}, queryParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks":           tasks,
		"pagination_data": paginationData,
	})
}
