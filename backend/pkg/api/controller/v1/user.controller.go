package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"github.com/whyxn/easynas/backend/pkg/dto"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/util"
	"net/http"
)

type UserControllerInterface interface {
	Create(c *gin.Context)
	GetList(c *gin.Context)
	Get(c *gin.Context)
	Delete(c *gin.Context)
}

type userController struct{}

var uc userController

func UserController() *userController {
	return &uc
}

// Create User
func (ctrl *userController) Create(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	var input dto.CreateUserInputDTO

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	user := &model.User{
		Name:        input.Name,
		Email:       input.Email,
		Password:    "",
		NasClientIP: input.NasClientIP,
		Role:        input.Role,
	}

	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		log.Logger.Errorw("Failed to hash password", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	user.Password = hashedPassword

	if err = db.GetDb().Insert(user); err != nil {
		log.Logger.Fatalw("Failed to create user", "err", err.Error())
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetList of User
func (ctrl *userController) GetList(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	userList, err := db.GetList[model.User](db.GetDb(), map[string]interface{}{})
	if err != nil {
		log.Logger.Errorw("Failed to fetch user list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   userList,
	})
}

// Get User
func (ctrl *userController) Get(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	id := ctx.Param("id")

	user, err := db.Get[model.User](db.GetDb(), map[string]interface{}{"ID": id})
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Logger.Errorw("Failed to fetch user list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// Delete User
func (ctrl *userController) Delete(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	id := ctx.Param("id")

	if err := db.GetDb().Delete(&model.User{}, map[string]interface{}{"ID": id}); err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
