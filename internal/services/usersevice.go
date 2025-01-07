package services

import "userManagement/internal/repository"

type UserService struct {
	userRepo repository.IUserRepo
}

type IUserService interface{

}

func NewUserService(userRepo repository.IUserRepo) *UserService{
	return &UserService{userRepo: userRepo}
}