package authorization

type Permission uint8

const (
	PermissionNone Permission = iota
	PermissionRead
	PermissionWrite
	PermissionDelete
)
