package page

import (
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func History() []byte {
	page := embedIntoApp(buildHistory(), "history", 4)
	return []byte(page.Body)
}

func buildHistory() *template.Page {
	const title = "История просмотров"
	const name = "history"
	args := struct {
		Title string
	}{
		Title: title,
	}

	builder := newPageBuilder()
	page, err := builder.Name(name).Template(html.HistoryTemplate).Args(args).Build(args.Title)

	if err != nil {
		panic(err)
	}

	return page
}
