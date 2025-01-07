package controllers

import "userManagement/internal/services"

type UserController struct {
	userService services.IUserService
}

type IUserController interface{

}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{userService: userService}
}