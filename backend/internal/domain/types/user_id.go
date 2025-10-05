package types

type UserID string

func (u UserID) String() string {
	return string(u)
}

func ParseUserID(s string) (UserID, error) {
	return UserID(s), nil
}

func ParseUserIDOrNil(s *string) (*UserID, error) {
	if s == nil {
		return nil, nil
	}

	id, err := ParseUserID(*s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
