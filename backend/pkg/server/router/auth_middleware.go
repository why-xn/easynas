package router

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/log"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if accessToken := c.GetHeader("Authorization"); len(accessToken) > 0 {
			//context.AddAccessTokenToContext(c, accessToken)
			//context.AddRequesterToContext(c, &requester)
		} else {
			// Access Token not found in request header
			log.Logger.Debug("Access Token not found in request header")
			// c.JSON(http.StatusUnauthorized, logError("invalid token"))
			// c.Abort()
		}
		c.Next()
	}
}

func logError(errMsg string) gin.H {
	return gin.H{"msg": errMsg}
}
