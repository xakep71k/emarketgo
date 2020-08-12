package html

import (
	"bytes"
	"fmt"
	"text/template"
)

func Generate(name, templateStr string, input interface{}) *bytes.Buffer {
	templ, err := template.New(name).Funcs(defaultTemplateFuncs()).Parse(templateStr)
	if err != nil {
		panic(fmt.Sprintf("%v %v", name, err))
	}

	buf := &bytes.Buffer{}
	if err := templ.Execute(buf, input); err != nil {
		panic(fmt.Sprintf("%v %v", name, err))
	}

	return buf
}

func defaultTemplateFuncs() map[string]interface{} {
	return template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
		"keyHistory": func() string {
			return "emarket.history.v1"
		},
		"keyCart": func() string {
			return "emarket.cart.v1"
		},
	}
}
