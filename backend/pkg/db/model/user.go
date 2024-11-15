package model

const (
	RoleAdmin = "ROLE_ADMIN"
	RoleUser  = "ROLE_USER"
)

type User struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"-"`
	NasClientIP string `json:"nasClientIP" gorm:"unique"`
	Role        string `json:"role"`
}
