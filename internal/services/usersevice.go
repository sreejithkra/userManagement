package services

import (
	"errors"
	"userManagement/internal/models"
	"userManagement/internal/repository"
)

type UserService struct {
	userRepo repository.IUserRepo
}

type IUserService interface {
	Signup(user *models.User) error

}

func NewUserService(userRepo repository.IUserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Signup(user *models.User) error{
	ExistingUser,_:=s.userRepo.GetUserByEmail(user.Email)
	if ExistingUser!=nil{
		return errors.New(models.UserExist)
	}

	err:=s.userRepo.CreateUser(user)
	if err!=nil{
		return err
	}
	return nil
}
