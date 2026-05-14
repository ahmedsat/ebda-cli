package utils

import (
	"strings"
	"unicode"
)

// SameAfterSanitize returns true if a and b are equal after full sanitization.
func SameAfterSanitize(a, b string) bool {
	return sanitize(a) == sanitize(b)
}

func sanitize(s string) string {
	var out []rune
	for _, r := range s {
		// drop non-printable
		if !unicode.IsPrint(r) {
			continue
		}
		// drop all whitespace (ASCII + Unicode)
		if unicode.IsSpace(r) {
			continue
		}

		// normalize Arabic variants
		switch r {

		// Alef variants → ا
		case 'أ', 'إ', 'آ', 'ٱ':
			r = 'ا'

		// Yeh variants → ي
		case 'ى', 'ئ':
			r = 'ي'

		// Waw with hamza → و
		case 'ؤ':
			r = 'و'

		// Teh Marbuta → ه
		case 'ة':
			r = 'ه'

		// Remove Arabic diacritics (tashkeel)
		case
			'َ', // fatha
			'ً', // tanween fatha
			'ُ', // damma
			'ٌ', // tanween damma
			'ِ', // kasra
			'ٍ', // tanween kasra
			'ْ', // sukoon
			'ّ': // shadda
			continue
		}

		// normalize case for non-Arabic scripts
		r = unicode.ToLower(r)

		out = append(out, r)
	}

	return strings.TrimSpace(string(out))
}

func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	var token []rune
	var capitalizeNext = true

	flushToken := func() {
		if len(token) == 0 {
			return
		}

		result = append(result, unicode.ToUpper(token[0]))
		if isAllUpper(token[1:]) {
			for _, r := range token[1:] {
				result = append(result, unicode.ToLower(r))
			}
		} else {
			result = append(result, token[1:]...)
		}

		token = token[:0]
	}

	for _, r := range s {
		switch {
		case r == '_' || r == '-' || r == ' ' || r == ':' || r == '(' || r == ')' || r == '/':
			flushToken()
			capitalizeNext = true
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if capitalizeNext {
				capitalizeNext = false
			}
			token = append(token, r)
		default:
			flushToken()
			// skip other symbols
			capitalizeNext = true
		}
	}
	flushToken()

	return string(result)
}

func isAllUpper(rs []rune) bool {
	for _, r := range rs {
		if unicode.IsLetter(r) && unicode.IsLower(r) {
			return false
		}
	}
	return true
}
