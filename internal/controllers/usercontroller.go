package controllers

import (
	"fmt"
	"net/http"
	"userManagement/internal/models"
	"userManagement/internal/services"
	"userManagement/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.IUserService
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
			ctx.JSON(http.StatusConflict, gin.H{"error": models.UserExist})
			return
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": models.SignupSuccess,
	})

}

func (c *UserController) Login(ctx *gin.Context) {
	var loginRequest models.UserLogin
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.InvalidInput})
		return
	}

	err := utils.Validate(loginRequest)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Login(&loginRequest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": models.UserNotFound})
		return
	}
	fmt.Println(loginRequest, user.Password)

	check := c.userService.VerifyPassword(&loginRequest, user)
	if !check {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.IncorrectPassword})
		return
	}

	accessToken, err := utils.GenerateJWT(user.Email, user.ID, "user", 1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Login Succesful", "token": accessToken})

}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userId := GetUserId(ctx)
	user, err := c.userService.GetProfile(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user data"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func GetUserId(ctx *gin.Context) string {
	claims, exists := ctx.Get("ID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Claims not found"})
		return ""
	}
	// Attempt to assert claims as float64
	userIDFloat, ok := claims.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return ""
	}
	// Convert the float64 to a string
	userID := fmt.Sprintf("%.0f", userIDFloat)
	return userID
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userId := GetUserId(ctx)

	var request models.UpdateUser
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.InvalidInput})
		return
	}

	err := utils.Validate(request)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = c.userService.UpdateProfile(userId, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})

}
