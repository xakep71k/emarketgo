package impl

import (
	"bytes"
	"emarket/html"
)

const (
	pageNo = iota
	pageHome
	pageDelivery
	pageContact
	pageHistory
	pageNewOrder
)

type Page struct {
	ID          string
	Title       string
	Body        string
	CurrentPage int
	PageTitle   string
}

type ProductPageList struct {
	First       bool
	Last        bool
	PageNum     int
	Products    []*Product
	PageNumbers []int
	MaxPages    int
	ListHTML    string
	Title       string
}

func (p *Page) htmlData() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	html.AppPage.Execute(buf, p)
	return doMinify(buf.Bytes(), "text/html")
}

func NewHomePage(body string) *Page {
	p := &Page{
		ID:          "magazines-index",
		Title:       ProductsTitle,
		CurrentPage: pageHome,
		Body:        body,
	}

	return p
}

func NewDeliveryPage() *Page {
	p := &Page{
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

func NewProductPageList(index int, maxPages int, products []*Product, listHtml string) *ProductPageList {
	p := &ProductPageList{
		PageNum:  index + 1,
		First:    index == 0,
		Last:     maxPages == 0 || maxPages == index+1,
		Products: products,
		MaxPages: maxPages,
		ListHTML: listHtml,
		Title:    ProductsTitle,
	}

	min := func(x, y int) int {
		if x > y {
			return y
		}
		return x
	}

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
