package dto

import "github.com/whyxn/easynas/backend/pkg/enum"

type LoginInputDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserInputDTO struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	NasClientIP     string `json:"nasClientIP"`
	Role            string `json:"role"`
}

type UpdateUserPasswordInputDTO struct {
	Id              uint   `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type DeleteUserInputDTO struct {
	Id uint `json:"id"`
}

type CreateZfsDatasetInputDTO struct {
	Pool        string `json:"pool"`
	DatasetName string `json:"datasetName"`
	Quota       string `json:"quota"`
}

type DeleteZfsDatasetInputDTO struct {
	Pool        string `json:"pool"`
	DatasetName string `json:"datasetName"`
}

type GetZfsDatasetInputDTO struct {
	Name string `json:"name"`
}

type CreateNfsShareInputDTO struct {
	Pool        string `json:"pool"`
	DatasetName string `json:"datasetName"`
}

type DeleteNfsShareInputDTO struct {
	Pool        string `json:"pool"`
	DatasetName string `json:"datasetName"`
}

type AddUserPermissionToNfsShareInputDTO struct {
	UserId      uint                `json:"userId"`
	DatasetName string              `json:"datasetName"`
	Permission  enum.PermissionType `json:"permission"`
}

type RemoveUserPermissionFromNfsShareInputDTO struct {
	Id uint `json:"id"`
}

type RestoreFromSnapshotInputDTO struct {
	SnapshotName string `json:"snapshotName"`
}
