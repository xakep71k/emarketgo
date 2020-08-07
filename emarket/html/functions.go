package html

import (
	"bytes"
	"fmt"
	"text/template"
)

func Generate(name, templateStr string, input interface{}) string {
	funcMap := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	}

	templ, err := template.New(name).Funcs(funcMap).Parse(templateStr)
	if err != nil {
		panic(fmt.Sprintf("%v %v", name, err))
	}

	var buf bytes.Buffer
	if err := templ.Execute(&buf, input); err != nil {
		panic(fmt.Sprintf("%v %v", name, err))
	}

	return buf.String()
}
