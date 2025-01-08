package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    FirstName string `json:"first_name" validate:"required,nameOrInitials"`
    LastName  string `json:"last_name" validate:"required,nameOrInitials"`
    Email     string `json:"email" gorm:"unique" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8,max=16,password"`
    Phone     string `json:"phone" gorm:"unique" validate:"required,numeric,len=10"`
    Status    string `gorm:"type:varchar(10);check:status IN ('Active','Blocked','Deleted');default:'Active'" json:"status" `
}

type UserLogin struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required"`
}

type UpdateUser struct {
    FirstName string `json:"first_name" validate:"required,nameOrInitials"`
    LastName  string `json:"last_name" validate:"required,nameOrInitials"`
    Phone     string `json:"phone" gorm:"unique" validate:"required,numeric,len=10"`
}


