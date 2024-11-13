package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string
	Email       string
	Password    string
	NasClientIP string `gorm:"unique"`
}
