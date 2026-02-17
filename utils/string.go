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
