package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string
	Email       string `gorm:"unique"`
	Password    string
	NasClientIP string `gorm:"unique"`
}
