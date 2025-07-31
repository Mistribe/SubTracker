package x

import "regexp"

func MakeKey(input string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	processedString := reg.ReplaceAllString(input, "-")
	return processedString
}
