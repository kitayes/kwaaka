package v1

import (
	"kwaaka/api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateProductStatus(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id is required"})
		return
	}

	var req models.ProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Info("UpdateProductStatus: invalid body: ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	// пока без JWT просто ставим "system" или "admin" на всякий
	userID := "system"

	if err := h.services.Product.ChangeStatus(
		c.Request.Context(),
		productID,
		req.Status,
		req.Reason,
		userID,
	); err != nil {
		h.logger.Error("UpdateProductStatus err: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to queue status change",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "status update queued",
	})
}
