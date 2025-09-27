package family

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/mistribe/subtracker/internal/domain/types"
)

func NewGenerateInvitationCode(familyID types.FamilyID) (string, error) {
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", nil
	}

	combined := familyID[:]
	combined = append(combined, randomBytes...)

	return base64.URLEncoding.EncodeToString(combined), nil
}
