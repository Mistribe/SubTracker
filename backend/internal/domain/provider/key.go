package provider

import (
	"strings"

	"github.com/mistribe/subtracker/internal/domain/types"
)

func sanitizeKey(input string) string {
	lowered := strings.ToLower(input)

	var b strings.Builder
	for _, r := range lowered {
		switch {
		case r == ' ':
			b.WriteByte('-')
		case r == '-':
			b.WriteByte('-')
		case r == '+':
			b.WriteString("-plus")
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			// all other runes are skipped
		}
	}

	return b.String()
}

func generateKey(name string, owner types.Owner) string {
	switch owner.Type() {
	case types.PersonalOwnerType:
		return "c_" + sanitizeKey(name) + "_" + owner.UserId().String()
	case types.FamilyOwnerType:
		return "f_" + sanitizeKey(name) + "_" + owner.FamilyId().String()
	case types.SystemOwnerType:
		return "s_" + sanitizeKey(name)
	default:
		panic("unknown owner type")
	}
}
