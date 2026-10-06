package wiringmap

import (
	"fmt"
	"strings"
	"unicode"
)

// SanitizeID maps paths and labels to wiringmap-safe id segments (S in the slice contract).
// '/' becomes '.'; any other rune outside [A-Za-z0-9._-] becomes '_'.
func SanitizeID(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '/':
			b.WriteByte('.')
		case r == '.' || r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

// CheckSanitizeCollisions returns an error if two distinct inputs collapse to the same S().
func CheckSanitizeCollisions(keys map[string]string) error {
	seen := make(map[string]string)
	for original, sanitized := range keys {
		if prev, ok := seen[sanitized]; ok {
			return fmt.Errorf("sanitize collision: %q and %q both map to %q", prev, original, sanitized)
		}
		seen[sanitized] = original
	}
	return nil
}
