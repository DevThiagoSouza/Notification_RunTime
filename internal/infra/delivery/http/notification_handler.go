package http

import (
	"apiRest/internal/entity"
	"apiRest/internal/infra/observability"
	"apiRest/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	UseCase *usecase.SendNotificationUseCase
	Metrics *observability.Metrics
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var dto struct {
		UserID  string `json:"user_id"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payload inválido: " + err.Error()})
		return
	}

	userID := dto.UserID
	if userID == "" {
		userID = c.GetString("user_id")
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	notification := entity.Notification{
		ID:      uuid.NewString(),
		UserID:  userID,
		Content: dto.Content,
		Type:    dto.Type,
		Status:  entity.StatusPending,
	}

	if err := h.UseCase.Execute(c.Request.Context(), &notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if h.Metrics != nil {
		h.Metrics.IncCreated()
	}

	c.JSON(http.StatusCreated, gin.H{"id": notification.ID, "status": notification.Status})
}
