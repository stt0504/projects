//go:build !solution

package ciletters

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed tmpl.txt
var text string

func lines10(s string) string {
	result := ""
	lines := strings.Split(s, "\n")
	var from int
	if len(lines) < 10 {
		from = 0
	} else {
		from = len(lines) - 10
	}
	for i := from; i < len(lines)-1; i++ {
		result += lines[i] + "\n            "
	}
	result += lines[len(lines)-1]
	return result
}

func MakeLetter(n *Notification) (string, error) {
	funcMap := template.FuncMap{
		"lines10": lines10,
	}

	tmpl := template.New("template").Funcs(funcMap)
	var err error
	tmpl, err = tmpl.Parse(text)
	if err != nil {
		return "", err
	}

	sb := new(strings.Builder)
	err = tmpl.Execute(sb, n)
	if err != nil {
		return "", err
	}

	result := sb.String()
	result = strings.ReplaceAll(result, "\r\n", "\n")

	return result, nil
}
