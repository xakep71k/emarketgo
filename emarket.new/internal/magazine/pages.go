package magazine

import (
	"emarket/internal/records"
	"errors"
)

const (
	pageNo = iota
	pageHome
	pageDelivery
	pageContact
	pageHistory
	pageNewOrder
)

type page struct {
	ID          string
	Title       string
	Body        string
	CurrentPage int
	PageTitle   string
}

type pagesList struct {
	First       bool
	Last        bool
	PageNum     int
	Magazines   []*records.Magazine
	PageNumbers []int
	MaxPages    int
	ListHTML    string
	Title       string
}

type pageWithMagazines struct {
	number    int
	magazines []*records.Magazine
}

const magazinesTitle = "Журналы и выкройки для шитья"

func newContactPage() (*page, error) {
	p := &page{
		Title:       "Контакты",
		CurrentPage: pageContact,
	}

	buf, err := execute("contact page", contactHTML, p)
	if err != nil {
		return nil, err
	}

	p.Body = buf.String()
	return p, nil
}

func newPageList(maxPages int, mpage pageWithMagazines, listHtml string) *pagesList {
	p := &pagesList{
		PageNum:   mpage.number,
		First:     mpage.number == 1,
		Last:      maxPages == 0 || maxPages == mpage.number,
		Magazines: mpage.magazines,
		MaxPages:  maxPages,
		ListHTML:  listHtml,
		Title:     magazinesTitle,
	}

	min := func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}

	index := mpage.number - 1
	pageNumbers := []int{index}
	prepend := true
	iPre := index - 1
	iApp := index + 1
	for len(pageNumbers) < min(maxPages, 5) {
		if prepend && iPre >= 0 {
			pageNumbers = append([]int{iPre}, pageNumbers...)
			iPre--
		} else if iApp < maxPages {
			pageNumbers = append(pageNumbers, iApp)
			iApp++
		}
		prepend = !prepend
	}

	p.PageNumbers = pageNumbers
	return p
}

func newHomePage(body string) *page {
	p := &page{
		ID:          "magazines-index",
		Title:       magazinesTitle,
		CurrentPage: pageHome,
		Body:        body,
	}

	return p
}

/*
func newDeliveryPage() *page {
	p := &page{
		ID:          "service-delivery_terms",
		Title:       "Условия доставки",
		CurrentPage: pageDelivery,
	}

	p.Body = html.Generate("delivery page", html.Delivery, p).String()
	return p
}

func NewContactPage() *Page {
	p := &Page{
		Title:       "Контакты",
		CurrentPage: pageContact,
	}

	p.Body = html.Generate("contact page", html.Contact, p).String()
	return p
}

func NewHistoryPage() *Page {
	p := &Page{
		ID:          "service-history",
		Title:       "История просмотров",
		CurrentPage: pageHistory,
	}

	p.Body = html.Generate("history page", html.History, p).String()
	return p
}

func NewNotFoundPage() *Page {
	return &Page{
		Title: "Страница не найдена",
		Body:  html.NotFound,
	}
}

func NewProductsPage(body string) *Page {
	return &Page{
		ID:          "magazines-index",
		Title:       "Журналы и выкройки для шитья",
		CurrentPage: pageHome,
		Body:        body,
	}
}

func NewProductPage(body, title string) *Page {
	return &Page{
		ID:          "magazines-show",
		Title:       title,
		CurrentPage: pageNo,
		Body:        body,
	}
}

func NewOrderPage() *Page {
	p := &Page{
		ID:          "orders-new",
		Title:       "Корзина",
		CurrentPage: pageNewOrder,
	}

	p.Body = html.Generate("cart page", html.NewOrder, p).String()
	return p
}
*/

func arrangeMagazinesPerPage(magazines []*records.Magazine) (magazinesPages []pageWithMagazines) {
	const pageSize = 30
	iPage := -1
	for i, p := range magazines {
		if (i % pageSize) == 0 {
			iPage++
			magazinesPages = append(magazinesPages, pageWithMagazines{})
			magazinesPages[iPage].number = iPage + 1
		}

		magazinesPages[iPage].magazines = append(magazinesPages[iPage].magazines, p)
	}

	return
}

func prepareHTMLMagazinesPages(magazinePages []pageWithMagazines) ([]string, error) {
	pages := make([]string, 0)
	maxPages := len(magazinePages)

	if maxPages == 0 {
		return nil, errors.New("there are no any magazines")
	}

	for _, mpage := range magazinePages {
		plist, err := execute("product list", magazineListHTML, mpage)
		if err != nil {
			return nil, err
		}

		pagination, err := execute(
			"pagination",
			paginationHTML,
			newPageList(maxPages, mpage, plist.String()),
		)

		if err != nil {
			return nil, err
		}

		pages = append(pages, pagination.String())
	}
	return pages, nil
}
