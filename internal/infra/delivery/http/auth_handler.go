package http

import (
	"net/http"

	"apiRest/internal/infra/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	JWTManager *auth.JWTManager
}

func (h *AuthHandler) Token(c *gin.Context) {
	var payload struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload: " + err.Error()})
		return
	}

	token, err := h.JWTManager.Generate(payload.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
