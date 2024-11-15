package model

type NfsShare struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	Pool    string `json:"pool"`
	Dataset string `json:"dataset"`
}

type NfsSharePermission struct {
	ID         uint     `json:"id" gorm:"primarykey"`
	NfsShareId uint     `json:"nfsShareId"`
	NfsShare   NfsShare `json:"nfsShare" gorm:"foreignKey:NfsShareId"`
	UserId     uint     `json:"userId"`
	User       User     `json:"user" gorm:"foreignKey:UserId"`
	Permission string   `json:"permission"`
}
