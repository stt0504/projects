//go:build !solution

package speller

import (
	"strings"
)

var m = map[int64]string{
	0:          "zero",
	1:          "one",
	2:          "two",
	3:          "three",
	4:          "four",
	5:          "five",
	6:          "six",
	7:          "seven",
	8:          "eight",
	9:          "nine",
	10:         "ten",
	11:         "eleven",
	12:         "twelve",
	13:         "thirteen",
	14:         "fourteen",
	15:         "fifteen",
	16:         "sixteen",
	17:         "seventeen",
	18:         "eighteen",
	19:         "nineteen",
	20:         "twenty",
	30:         "thirty",
	40:         "forty",
	50:         "fifty",
	60:         "sixty",
	70:         "seventy",
	80:         "eighty",
	90:         "ninety",
	100:        "hundred",
	1000:       "thousand",
	1000000:    "million",
	1000000000: "billion",
}

func Spell(n int64) string {
	if n == 0 {
		return m[0]
	}
	var builder strings.Builder
	builder.Grow(500)
	if n < 0 {
		builder.WriteString("minus ")
		n = -n
	}

	spellUtil(n, &builder)

	return strings.TrimSpace(builder.String())
}

func spellUtil(num int64, builder *strings.Builder) string {
	if num >= 1000000000 {
		spellUtil(num/1000000000, builder)
		builder.WriteString(m[1000000000] + " ")
		num %= 1000000000
	}
	if num >= 1000000 {
		spellUtil(num/1000000, builder)
		builder.WriteString(m[1000000] + " ")
		num %= 1000000
	}
	if num >= 1000 {
		spellUtil(num/1000, builder)
		builder.WriteString(m[1000] + " ")
		num %= 1000
	}
	if num >= 100 {
		builder.WriteString(m[num/100] + " " + m[100] + " ")
		num %= 100
	}
	if num >= 20 {
		builder.WriteString(m[num-num%10])
		num %= 10
		if num > 0 {
			builder.WriteString("-")
		} else {
			builder.WriteString(" ")
		}
	}
	if num > 0 {
		builder.WriteString(m[num] + " ")
	}
	return builder.String()
}
