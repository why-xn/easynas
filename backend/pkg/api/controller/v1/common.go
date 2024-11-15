package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/db/model"
)

func returnErrorResponse(ctx *gin.Context, msg string, statusCode int) {
	ctx.JSON(statusCode, gin.H{
		"status": "error",
		"msg":    msg,
	})
}

func isAdmin(requester *model.User) bool {
	if requester.Role == model.RoleAdmin {
		return true
	}
	return false
}
