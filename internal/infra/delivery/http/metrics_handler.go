package http

import (
	"net/http"

	"apiRest/internal/infra/observability"
	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	Metrics *observability.Metrics
}

func (h *MetricsHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, h.Metrics.Snapshot())
}
