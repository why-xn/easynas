package model

import "gorm.io/gorm"

type NfsShare struct {
	gorm.Model
	Pool    string
	Dataset string
}

type NfsSharePermission struct {
	gorm.Model
	NfsShareId uint
	NfsShare   NfsShare `gorm:"foreignKey:NfsShareId"`
	UserId     uint
	User       User `gorm:"foreignKey:UserId"`
	Permission string
}
