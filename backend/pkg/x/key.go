package x

import (
	"strings"
)

func MakeKey(input string) string {
	input = strings.ToLower(input)
	input = strings.ReplaceAll(input, " ", "-")
	return input
}
