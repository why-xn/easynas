package model

import "gorm.io/gorm"

const (
	RoleAdmin = "ROLE_ADMIN"
	RoleUser  = "ROLE_USER"
)

type User struct {
	gorm.Model
	Name        string
	Email       string `gorm:"unique"`
	Password    string
	NasClientIP string `gorm:"unique"`
	Role        string
}
