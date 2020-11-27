package page

import (
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func NotFound() []byte {
	page := embedIntoApp(buildNotFound(), "notfound", 0)
	return []byte(page.Body)
}

func buildNotFound() *template.Page {
	const title = "страница не найдена"
	const name = "notfound"
	args := struct {
		Title string
	}{
		Title: title,
	}

	builder := newPageBuilder()
	page, err := builder.Name(name).Template(html.NotFound).Args(args).Build(args.Title)

	if err != nil {
		panic(err)
	}

	return page
}
