package main

import (
	"fmt"
	"userManagement/internal/config"
	"userManagement/internal/controllers"
	"userManagement/internal/database"
	"userManagement/internal/middleware"
	"userManagement/internal/repository"
	"userManagement/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	env := config.EnvConfig()
	db := database.ConnectDatabase(env)

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	fmt.Println(userController)

	router.POST("signup", userController.Signup)
	router.POST("login", userController.Login)

	userGroup := router.Group("user")
	userGroup.Use(middleware.JWTMIddleware("user"))
	userGroup.GET("profile", userController.GetProfile)
	userGroup.PUT("profile", userController.UpdateProfile)

	err := router.Run(":8086")
	if err != nil {
		fmt.Println("failed to start server")
	}

}
