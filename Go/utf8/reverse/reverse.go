//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	for len(input) > 0 {
		r, size := utf8.DecodeLastRuneInString(input)
		input = input[:len(input)-size]
		builder.WriteRune(r)
	}
	return builder.String()
}
