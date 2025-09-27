package types

type UserID string

func (u UserID) String() string {
	return string(u)
}

func ParseUserID(s string) (UserID, error) {
	return UserID(s), nil
}
