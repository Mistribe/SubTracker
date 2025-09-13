package user

type Permission string

const (
	PermissionRead   Permission = "READ"
	PermissionWrite  Permission = "WRITE"
	PermissionDelete Permission = "DELETE"
)
