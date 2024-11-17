package enum

type PermissionType string

const (
	ReadOnly  PermissionType = "r"
	ReadWrite PermissionType = "rw"
)
