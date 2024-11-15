package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/whyxn/easynas/backend/pkg/api/controller/v1"
)

func AddApiRoutes(httpRg *gin.RouterGroup) {
	httpRg.GET("health", v1.HealthController().Check)
	httpRg.GET("health/secured", v1.HealthController().SecuredCheck)

	httpRg.POST("api/v1/auth/login", v1.AuthController().Login)

	httpRg.POST("api/v1/users", v1.UserController().Create)
	httpRg.GET("api/v1/users/:id", v1.UserController().Get)
	httpRg.GET("api/v1/users", v1.UserController().GetList)
	httpRg.DELETE("api/v1/users/:id", v1.UserController().Delete)
}
