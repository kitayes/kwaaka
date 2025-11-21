package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetMenu(c *gin.Context) {
	menuID := c.Param("menu_id")
	if menuID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		return
	}

	menu, err := h.services.Menu.GetMenu(c.Request.Context(), menuID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "menu not found"})
		return
	}

	c.JSON(http.StatusOK, menu)
}
