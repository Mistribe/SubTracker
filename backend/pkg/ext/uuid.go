package ext

import "github.com/google/uuid"

func UuidToString(source uuid.UUID) string {
	return source.String()
}
