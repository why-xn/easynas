package model

import "github.com/whyxn/easynas/backend/pkg/enum"

type NfsShare struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Pool    string `json:"pool"`
	Dataset string `json:"dataset"`
	ShareOn bool   `json:"shareOn"`
}

type NfsSharePermission struct {
	ID         uint                `json:"id" gorm:"primarykey"`
	NfsShareId uint                `json:"-"`
	NfsShare   NfsShare            `json:"nfsShare" gorm:"foreignKey:NfsShareId"`
	UserId     uint                `json:"-"`
	User       User                `json:"user" gorm:"foreignKey:UserId"`
	Permission enum.PermissionType `json:"permission"`
}
