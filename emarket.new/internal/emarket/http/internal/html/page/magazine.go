package page

import (
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func Magazine(magaz *emarket.Magazine) []byte {
	page := embedIntoApp(buildMagazine(magaz), "show magazine", 0)
	return []byte(page.Body)
}

func buildMagazine(magaz *emarket.Magazine) *template.Page {
	const name = "show magazine"

	builder := newPageBuilder()
	page, err := builder.Name(name).Template(html.MagazineTemplate).Args(magaz).Build(magaz.Title)

	if err != nil {
		panic(err)
	}

	return page
}
