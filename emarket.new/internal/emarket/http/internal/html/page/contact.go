package page

import (
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func Contact() []byte {
	page := embedIntoApp(buildContact(), "contact", 3)
	return []byte(page.Body)
}

func buildContact() *template.Page {
	const title = "Контакты"
	const name = "contact"
	args := struct {
		Title string
	}{
		Title: title,
	}

	builder := newPageBuilder()
	page, err := builder.Name(name).Template(html.ContactTemplate).Args(args).Build(args.Title)

	if err != nil {
		panic(err)
	}

	return page
}
