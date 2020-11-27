package page

import (
	"emarket/internal/emarket/http/internal/html"
	"emarket/internal/pkg/template"
)

func embedIntoApp(page *template.Page, id string, curPageNum int) *template.Page {
	args := struct {
		Title          string
		CurrentPageNum int
		Body           string
		ID             string
	}{
		Title:          page.Title,
		CurrentPageNum: curPageNum,
		Body:           page.Body,
		ID:             id,
	}
	const name = "embedded into app page"

	builder := newPageBuilder()
	appPage, err := builder.Name(name).Template(html.AppTemplate).Args(args).Build(args.Title)

	if err != nil {
		panic(err)
	}

	return appPage
}

func embedIntoPagination(page *template.Page, index int, maxPages int) *template.Page {
	args := struct {
		Title       string
		PageNum     int
		ListHTML    string
		First       bool
		Last        bool
		MaxPages    int
		PageNumbers []int
	}{
		Title:       page.Title,
		PageNum:     index + 1,
		ListHTML:    page.Body,
		First:       index == 0,
		Last:        maxPages == 0 || maxPages == index+1,
		MaxPages:    maxPages,
		PageNumbers: genPageSlice(index, maxPages),
	}
	const name = "embedded into pagination"

	builder := newPageBuilder()
	page, err := builder.Name(name).Template(html.PaginationTemplate).Args(args).Build(args.Title)

	if err != nil {
		panic(err)
	}

	return page
}

func genPageSlice(index, maxPages int) []int {
	pages := []int{index}
	prepend := true
	iPre := index - 1
	iApp := index + 1
	min := func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}

	for len(pages) < min(maxPages, 5) {
		if prepend && iPre >= 0 {
			pages = append([]int{iPre}, pages...)
			iPre--
		} else if iApp < maxPages {
			pages = append(pages, iApp)
			iApp++
		}
		prepend = !prepend
	}

	return pages
}
