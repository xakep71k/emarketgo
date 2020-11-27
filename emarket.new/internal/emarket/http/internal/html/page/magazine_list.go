package page

import (
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

const magazinesOnPage = 30
const magazineTitle = "Журналы и выкройки для шитья"

func MagazineList(all []*emarket.Magazine) [][]byte {
	var rawPages [][]byte
	pages := buildMagazineList(all)

	for _, page := range pages {
		rawPages = append(rawPages, []byte(page.Body))
	}

	return rawPages
}

func buildMagazineList(all []*emarket.Magazine) []*template.Page {
	var pages []*template.Page
	magazinePages := arrangeMagazinesPerPage(all)
	builder := newPageBuilder()
	maxPages := len(magazinePages)

	for pageIndex, magazPage := range magazinePages {
		args := struct {
			Title string
		}{
			Title: magazineTitle,
		}

		magazPageHTML, err := builder.Name("magazine page").Template(html.MagazineListTemplate).Args(magazPage).Build(args.Title)

		if err != nil {
			panic(err)
		}

		pageWithPaginationHTML := embedIntoPagination(magazPageHTML, pageIndex, maxPages)
		appIndex := 0

		if pageIndex == 0 {
			appIndex = 1
		}

		page := embedIntoApp(pageWithPaginationHTML, "magazines-index", appIndex)

		pages = append(pages, page)
	}

	return pages
}

func arrangeMagazinesPerPage(magazines []*emarket.Magazine) (magazinePages [][]*emarket.Magazine) {
	iPage := -1

	for i, magaz := range magazines {
		if (i % magazinesOnPage) == 0 {
			iPage++
			magazinePages = append(magazinePages, make([]*emarket.Magazine, 0))
		}

		magazinePages[iPage] = append(magazinePages[iPage], magaz)
	}

	return
}
