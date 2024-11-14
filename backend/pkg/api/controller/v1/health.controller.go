package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"net/http"
)

type HealthControllerInterface interface {
	Check(c *gin.Context)
	SecuredCheck(c *gin.Context)
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

func (ctrl *healthController) SecuredCheck(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "unauthorized request",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Hi! %s", requester.Name),
	})
}
