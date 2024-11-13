package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthControllerInterface interface {
	Check(c *gin.Context)
}

type healthController struct{}

var hc healthController

func HealthController() *healthController {
	return &hc
}

func (ctrl *healthController) Check(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "i am alive",
	})
}
