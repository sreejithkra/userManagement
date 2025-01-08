package repository

import (
	"userManagement/internal/models"

	"gorm.io/gorm"
)

type UserRepo struct {
	Database *gorm.DB
}

type IUserRepo interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{Database: db}
}

func (repo *UserRepo) CreateUser(user *models.User) error {
	if err := repo.Database.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.Database.Where("email=?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &models.User{}, nil
}
