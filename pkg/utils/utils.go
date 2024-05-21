package utils

import "strings"

// Capitalize returns a capitalized string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
