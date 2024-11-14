package router

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/jwt"
	"github.com/whyxn/easynas/backend/pkg/log"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if accessToken := c.GetHeader("Authorization"); len(accessToken) > 0 {
			claims, err := jwt.ValidateJWT(accessToken)
			if err != nil {
				log.Logger.Warnw("Failed to validate JWT token", "err", err.Error())
			} else {
				context.AddAccessTokenToContext(c, accessToken)
				context.AddRequesterToContext(c, &claims.User)
			}
		} else {
			// Access Token not found in request header
			log.Logger.Debug("Access Token not found in request header")
			// c.JSON(http.StatusUnauthorized, logError("invalid token"))
			// c.Abort()
		}
		c.Next()
	}
}
