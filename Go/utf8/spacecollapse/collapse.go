//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)
		input = input[size:]
		if unicode.IsSpace(r) {
			builder.WriteRune(' ')
			for len(input) > 0 && unicode.IsSpace(r) {
				r, size = utf8.DecodeRuneInString(input)
				input = input[size:]
			}
		}
		if !unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
