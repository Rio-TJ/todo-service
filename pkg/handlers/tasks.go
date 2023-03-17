package handlers

import (
	"errors"
	"github.com/SoftclubIT/todo-service/pkg/helpers"
	"github.com/SoftclubIT/todo-service/pkg/models"
	"github.com/SoftclubIT/todo-service/pkg/scopes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// GetTasks godoc
//
//	@Summary		Get all the tasks
//	@Description	Get the list of all the tasks
//	@Tags			tasks
//	@Produce		json
//	@Param			page		query		int		false	"Pagination parameter: page number"
//	@Param			per_page	query		int		false	"Pagination parameter: number of elements on a page"
//	@Param			status		query		string	false	"Filtering by status (done|undone)"
//	@Success		200			{array}		models.Task
//	@Failure		500			{object}	models.HTTPError
//	@Router			/tasks [get]
func (h *Handler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	queryParams := c.Request.URL.Query()

	result := h.DB.
		Scopes(scopes.Paginate(queryParams), scopes.Filter(queryParams)).
		Find(&tasks)
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

// CreateTask godoc
//
//	@Summary		Create a task
//	@Description	Create a task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			text	body		string	true	"The text of the task"
//	@Success		200		{object}	models.Task
//	@Failure		400		{object}	models.HTTPError
//	@Failure		422		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/tasks [post]
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
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
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

// DeleteTask godoc
//
//	@Summary		Delete a task
//	@Description	Delete a task
//	@Tags			tasks
//	@Produce		json
//	@Param			taskID	path		int	true	"The ID of the task that should be deleted"
//	@Success		200		{string}	"Successfully deleted the task"
//	@Failure		404		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/tasks/{taskID} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	taskID := c.Param("taskID")

	var task models.Task

	result := h.DB.First(&task, taskID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Record not found",
		})
		return
	}

	h.DB.Delete(&task)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted task",
	})
}

// CompleteTask godoc
//
//	@Summary		Complete a task
//	@Description	Change status of a task to done
//	@Tags			tasks
//	@Produce		json
//	@Param			taskID	path		int	true	"The ID of the task which status should be changed to done"
//	@Success		200		{object}	models.Task
//	@Failure		404		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/tasks/{taskID}/completed [post]
func (h *Handler) CompleteTask(c *gin.Context) {
	taskID := c.Param("taskID")

	var task models.Task

	result := h.DB.First(&task, taskID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Record not found",
		})
		return
	}

	task.Status = "done"

	if result = h.DB.Save(&task); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UndoTask godoc
//
//	@Summary		Undo a task
//	@Description	Change status of a task to undone
//	@Tags			tasks
//	@Produce		json
//	@Param			taskID	path		int	true	"The ID of the task which status should be changed to undone"
//	@Success		200		{object}	models.Task
//	@Failure		404		{object}	models.HTTPError
//	@Failure		500		{object}	models.HTTPError
//	@Router			/tasks/{taskID}/completed [delete]
func (h *Handler) UndoTask(c *gin.Context) {
	taskID := c.Param("taskID")

	var task models.Task

	result := h.DB.First(&task, taskID)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Record not found",
		})
		return
	}

	task.Status = "undone"

	if result = h.DB.Save(&task); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, task)
}
