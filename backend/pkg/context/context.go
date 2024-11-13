package context

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/db/model"
)

func AddAccessTokenToContext(c *gin.Context, accessToken interface{}) {
	c.Set("AccessToken", accessToken)
}

func AddRequesterToContext(c *gin.Context, user *model.User) {
	c.Set("User", user)
}

func GetRequesterFromContext(c *gin.Context) *model.User {
	if val, ok := c.Get("User"); ok {
		if val == nil {
			return nil
		}
		user, ok := val.(*model.User)
		if !ok {
			return nil
		}
		return user
	}
	return nil
}
