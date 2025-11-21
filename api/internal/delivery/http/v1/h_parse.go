package v1

import (
	"kwaaka/api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateParseTask(c *gin.Context) {
	var req models.ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Info("c.ShouldBindJSON(...) err: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	taskID, err := h.services.Parse.CreateTask(
		c.Request.Context(),
		req.SpreadsheetID,
		req.RestaurantName,
	)
	if err != nil {
		h.logger.Error("CreateParseTask(...) err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create parsing task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"status":  "queued",
	})
}

func (h *Handler) GetParseTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	task, err := h.services.Parse.GetTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
