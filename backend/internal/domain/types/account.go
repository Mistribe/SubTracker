package types

import "errors"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) String() string {
	return string(r)
}

var (
	ErrUnknownRole = errors.New("unknown role")
)

func ParseRoleOrDefault(role string, defaultValue Role) Role {
	r, err := ParseRole(role)
	if err != nil {
		return defaultValue
	}
	return r
}

func ParseRole(role string) (Role, error) {
	switch role {
	case string(RoleAdmin):
		return RoleAdmin, nil
	case string(RoleUser):
		return RoleUser, nil
	}
	return "", ErrUnknownRole
}

func MustParseRole(role string) Role {
	r, err := ParseRole(role)
	if err != nil {
		panic(err)
	}
	return r
}
