package handlers

import (
	"errors"
	"github.com/SoftclubIT/todo-service/pkg/helpers"
	"github.com/SoftclubIT/todo-service/pkg/models"
	"github.com/SoftclubIT/todo-service/pkg/scopes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
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

func (h *Handler) CreateTask(c *gin.Context) {
	var task models.Task

	validationErrors := make(map[string]string)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Printf("%v\t%T", err, err)
		var valErrors validator.ValidationErrors
		if errors.As(err, &valErrors) {
			for _, valError := range valErrors {
				if _, exists := validationErrors[valError.Field()]; !exists {
					validationErrors[valError.Field()] = helpers.ValidationMessageForTag(valError.Tag(), valError.Param())
				}
			}
		} else if err == io.EOF {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Empty body",
			})
			return
		}
	}

	if len(validationErrors) != 0 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"errors": validationErrors,
		})
		return
	}

	task.Status = "undone"

	if result := h.DB.Save(&task); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}
