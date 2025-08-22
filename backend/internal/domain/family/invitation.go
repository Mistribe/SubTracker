package family

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/google/uuid"
)

func NewGenerateInvitationCode(familyId uuid.UUID) (string, error) {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", nil
	}

	combined := familyId[:]
	combined = append(combined, randomBytes...)

	return base64.URLEncoding.EncodeToString(combined), nil
}
