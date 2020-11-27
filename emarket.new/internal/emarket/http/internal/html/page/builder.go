package page

import (
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func newPageBuilder() *template.PageBuilder {
	pb := template.NewPageBuilder()
	return pb.TemplateFuncs(defaultTemplateFuncs())
}

func defaultTemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
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
			return html.AlertCartPutIn
		},
		"alertCartRemove": func() string {
			return html.AlertCartRemove
		},
		"iterate": func(start int, end int) (res []int) {
			for i := start; i <= end; i++ {
				res = append(res, i)
			}
			return
		},
	}
}
