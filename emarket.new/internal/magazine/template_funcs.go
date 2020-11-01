package magazine

import "text/template"

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
		"iterate": func(start int, end int) (res []int) {
			for i := start; i <= end; i++ {
				res = append(res, i)
			}
			return
		},
	}
}

const alertCartPutIn = `<div id="alertCart" class="alert alert-success alert-dismissible fade text-center z-depth-2" role="alert"><strong><a href="/zakazy/novyy">Перейти в корзину >>>&nbsp;</a><span id="alertPutInCartCounter" class="counter">1</span></strong><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>`

const alertCartRemove = `<div id="alertCart" class="alert alert-success alert-dismissible fade text-center z-depth-2" role="alert"><strong>Товар убран из <a href="/zakazy/novyy" class="underlineLink">корзины</a></strong><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button></div>`
