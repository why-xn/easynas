package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"github.com/whyxn/easynas/backend/pkg/dto"
	"github.com/whyxn/easynas/backend/pkg/jwt"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/util"
	"net/http"
)

type AuthControllerInterface interface {
	Login(c *gin.Context)
}

type authController struct{}

var ac authController

func AuthController() *authController {
	return &ac
}

func (ctrl *authController) Login(ctx *gin.Context) {
	var input dto.LoginInputDTO

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		return
	}

	user, _ := db.Get[model.User](db.GetDb(), map[string]interface{}{"email": input.Username})
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credential, user not found",
		})
		return
	}

	if !util.CheckPasswordHash(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid password",
		})
		return
	}

	authToken, err := jwt.GenerateJWT(*user)
	if err != nil {
		log.Logger.Errorw("Failed to generate JWT token", "err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": authToken,
	})
}
