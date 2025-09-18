package utils

import (
	"strings"
	"unicode/utf8"
)

// SanitizeText removes or replaces invalid UTF-8 characters and controls the length
func SanitizeText(text string, maxLength int) string {
	if text == "" {
		return ""
	}

	// Remove or replace invalid UTF-8 characters
	if !utf8.ValidString(text) {
		// Convert to valid UTF-8 by replacing invalid sequences
		text = strings.ToValidUTF8(text, "")
	}

	// Trim whitespace
	text = strings.TrimSpace(text)

	// Truncate if too long (accounting for UTF-8 characters)
	if len(text) > maxLength {
		// Find a safe truncation point that doesn't break UTF-8 characters
		if maxLength <= 0 {
			return ""
		}
		
		runes := []rune(text)
		if len(runes) > maxLength {
			text = string(runes[:maxLength])
		}
		
		// Double-check the byte length after rune truncation
		for len(text) > maxLength {
			runes := []rune(text)
			if len(runes) == 0 {
				return ""
			}
			text = string(runes[:len(runes)-1])
		}
	}

	return text
}