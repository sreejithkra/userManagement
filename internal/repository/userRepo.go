package repository

import "gorm.io/gorm"

type UserRepo struct {
	Database *gorm.DB
}

type IUserRepo interface {
}

func NewUserRepository(db *gorm.DB) *UserRepo{
	return &UserRepo{Database: db}
}
