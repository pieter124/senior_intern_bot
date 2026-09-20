package util

import (
	"strings"
	"unicode"
)

type Empty struct{}
type Set[T comparable] map[T]Empty

// ContainsAny is a pure function that reports whether text 's' mentions any of
// 'phrases' as whole words. Matching is case-insensitive and ignores
// punctuation, e.g. "Software Intern" matches "intern" but "International"
// does not. Both 's' and 'phrases' may be raw text.
func ContainsAny(s string, phrases []string) bool {
	// Lowercase, collapse non-alphanumeric runs into single spaces, and pad
	// with a space each side so a plain substring search matches whole words.
	normalise := func(s string) string {
		words := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})
		return " " + strings.Join(words, " ") + " "
	}

	s = normalise(s)
	for _, p := range phrases {
		if strings.Contains(s, normalise(p)) {
			return true
		}
	}
	return false
}
