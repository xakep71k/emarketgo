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
		"alertCartPutIn": func() string {
			return alertCartPutIn
		},
		"alertCartRemove": func() string {
			return alertCartRemove
		},
	}
}

var alertCartPutIn = `<div id="alertCart" class="alert alert-success alert-dismissible fade text-center z-depth-2" role="alert"><strong><a href="/zakazy/novyy">Перейти в корзину >>>&nbsp;</a><span id="alertPutInCartCounter" class="counter">1</span></strong><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>`

var alertCartRemove = `<div id="alertCart" class="alert alert-success alert-dismissible fade text-center z-depth-2" role="alert"><strong>Товар убран из <a href="/zakazy/novyy" class="underlineLink">корзины</a></strong><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>`
