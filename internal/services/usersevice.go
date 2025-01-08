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
	Login(user *models.UserLogin) (*models.User, error)
	VerifyPassword(request *models.UserLogin, user *models.User) bool
	GetProfile(userID string) (*models.User, error)
	UpdateProfile(userId string, user models.UpdateUser) error
}

func NewUserService(userRepo repository.IUserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Signup(user *models.User) error {
	ExistingUser, _ := s.userRepo.GetUserByEmail(user.Email)
	if ExistingUser != nil {
		return errors.New(models.UserExist)
	}

	err := s.userRepo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(loginRequest *models.UserLogin) (*models.User, error) {
	user, _ := s.userRepo.GetUserByEmail(loginRequest.Email)
	if user == nil {
		return nil, errors.New(models.UserNotFound)
	}
	return user, nil
}

func (s *UserService) VerifyPassword(request *models.UserLogin, user *models.User) bool {
	return user.Password == request.Password
}

func (c *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := c.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *UserService) UpdateProfile(userId string, user models.UpdateUser) error {
	User, err := c.userRepo.GetUserById(userId)
	if err != nil {
		return err
	}
	User.FirstName = user.FirstName
	User.LastName = user.LastName
	User.Phone = user.Phone
	err = c.userRepo.UpdateProfile(User)
	if err != nil {
		return err
	}
	return nil
}
