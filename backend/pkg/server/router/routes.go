package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/whyxn/easynas/backend/pkg/api/controller/v1"
)

func AddApiRoutes(httpRg *gin.RouterGroup) {
	httpRg.GET("health", v1.HealthController().Check)
}
