package varfmt

import (
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var builder strings.Builder
	counter := -1
	i := 0
	index := 0
	for i < len(format) {
		if format[i] == '{' {
			counter++
			builder.WriteString(format[index:i])
			i++
			number := 0
			hasNumber := false
			for format[i] != '}' {
				hasNumber = true
				number = number*10 + int(format[i]-'0')
				i++
			}
			index = i + 1
			var arg interface{}
			if hasNumber {
				arg = args[number]
			} else {
				arg = args[counter]
			}
			switch t := arg.(type) {
			case string:
				builder.WriteString(t)
			case int:
				builder.WriteString(strconv.Itoa(t))
			case float64:
				builder.WriteString(strconv.FormatFloat(t, 'g', -1, 64))
			case bool:
				builder.WriteString(strconv.FormatBool(t))
			}
		}
		i++
	}
	builder.WriteString(format[index:])
	return builder.String()
}
