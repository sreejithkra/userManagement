package controllers

import (
	"net/http"
	"userManagement/internal/models"
	"userManagement/internal/services"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.IUserService
}

type IUserController interface {
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) Signup(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.InvalidInput})
		return
	}

	err := utils.Validate(user)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = c.userService.Signup(&user)
	if err != nil {
		if err.Error() == models.UserExist {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": models.UserExist})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": models.SignupSuccess,
	})

}
