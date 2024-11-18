package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/metrics"
	"net/http"
)

type MetricsControllerInterface interface {
	GetSystemMetrics(c *gin.Context)
}

type metricsController struct{}

var mc metricsController

func MetricsController() *metricsController {
	return &mc
}

func (ctrl *metricsController) GetSystemMetrics(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	metrics, err := metrics.GetSystemMetrics()
	if err != nil {
		log.Logger.Errorw("failed to fetch system metrics", err.Error())
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   metrics,
	})
}
