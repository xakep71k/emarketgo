package impl

import (
	"bytes"
	"emarket/html"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"syscall"
)

const (
	PageNo = iota
	PageHome
	PageDelivery
	PageContact
	PageHistory
	PageNewOrder
)

const pageSize = 30

type Product struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	Title       string `bson:"title" json:"title"`
	Price       int    `bson:"price" json:"price"`
	Thumb       []byte `bson:"thumb" json:"thumb"`
	Enable      bool   `bson:"enable" json:"enable"`
	Description string `bson:"description" json:"description"`
	Quantity    int    `bson:"quantity" json:"quantity"`
	OldID       int    `bson:"oldid" json:"oldid"`
	OldImgName  string `bson:"oldimgfile" json:"oldimgfile"`
	PageNum     int    `bson:"-" json:"-"`
}

type Cart struct {
	Products   []*Product
	Empty      bool
	TotalPrice int
}

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

func (p *Page) HTMLData() []byte {
	buf := bytes.NewBuffer(make([]byte, 0))
	html.AppPage.Execute(buf, p)
	return buf.Bytes()
}

const ProductsTitle = "Журналы и выкройки для шитья"

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

func NewHomePage(body string) *Page {
	p := &Page{
		ID:          "magazines-index",
		Title:       ProductsTitle,
		CurrentPage: PageHome,
		Body:        body,
	}

	return p
}

func NewDeliveryPage() *Page {
	p := &Page{
		ID:          "service-delivery_terms",
		Title:       "Условия доставки",
		CurrentPage: PageDelivery,
	}

	p.Body = html.Generate("delivery page", html.Delivery, p).String()
	return p
}

func NewContactPage() *Page {
	p := &Page{
		Title:       "Контакты",
		CurrentPage: PageContact,
	}

	p.Body = html.Generate("contact page", html.Contact, p).String()
	return p
}

func NewHistoryPage() *Page {
	p := &Page{
		ID:          "service-history",
		Title:       "История просмотров",
		CurrentPage: PageHistory,
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
		CurrentPage: PageHome,
		Body:        body,
	}
}

func NewProductPage(body, title string) *Page {
	return &Page{
		ID:          "magazines-show",
		Title:       title,
		CurrentPage: PageNo,
		Body:        body,
	}
}

func NewOrderPage() *Page {
	p := &Page{
		ID:          "orders-new",
		Title:       "Корзина",
		CurrentPage: PageNewOrder,
	}

	p.Body = html.Generate("cart page", html.NewOrder, p).String()
	return p
}

type Content struct {
	mytype map[string]string
}

func NewContent() *Content {
	c := &Content{}
	c.mytype = map[string]string{
		".js":     "application/javascript",
		".css":    "text/css",
		".ico":    "image/x-icon",
		".gif":    "image/gif",
		"js.map":  "application/octet-stream",
		"css.map": "application/octet-stream",
		".woff":   "font/woff",
		".woff2":  "font/woff2",
		".ttf":    "font/ttf",
	}
	return c
}

type EMarket struct {
	rootDir     string
	content     *Content
	router      *http.ServeMux
	Pages       map[string][]byte
	ProductsMap map[string]*Product
}

func buildHtmlProductPages(productPages [][]*Product) []string {
	pages := make([]string, 0)
	maxPages := len(productPages)

	if maxPages == 0 {
		panic("there are no any products")
	}

	for index, products := range productPages {
		plist := html.Generate("product list", html.ProductList, products).String()
		pagination := html.Generate(
			"pagination",
			html.Pagination,
			NewProductPageList(index, maxPages, products, plist),
		).String()
		pages = append(pages, pagination)
	}
	return pages
}

func collectProductsPages(products []*Product) (productPages [][]*Product) {
	iPage := -1
	for i, p := range products {
		if (i % pageSize) == 0 {
			iPage++
			productPages = append(productPages, make([]*Product, 0))
		}

		productPages[iPage] = append(productPages[iPage], p)
	}

	return
}

func setPageNum(productPages [][]*Product) {
	for index, products := range productPages {
		for _, p := range products {
			p.PageNum = index + 1
		}
	}
}

func NewEMarket(rootDir string, products []*Product) (*EMarket, error) {
	sort.SliceStable(products, func(i, j int) bool {
		return products[i].Title < products[j].Title
	})

	e := &EMarket{
		rootDir:     rootDir,
		content:     NewContent(),
		ProductsMap: make(map[string]*Product, 0),
	}

	for _, product := range products {
		e.ProductsMap[product.ID] = product
	}

	productPages := collectProductsPages(products)
	setPageNum(productPages)
	productPagesHtml := buildHtmlProductPages(productPages)
	e.Pages = map[string][]byte{
		"contact":  NewContactPage().HTMLData(),
		"delivery": NewDeliveryPage().HTMLData(),
		"notfound": NewNotFoundPage().HTMLData(),
		"home":     NewHomePage(productPagesHtml[0]).HTMLData(),
		"history":  NewHistoryPage().HTMLData(),
		"neworder": NewOrderPage().HTMLData(),
	}

	e.setupRouter(products, productPagesHtml)
	return e, nil
}

func (e *EMarket) setupRouter(products []*Product, productPagesHtml []string) {
	router := http.NewServeMux()
	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		e.handleSpecifiedFile("/static/favicon.ico", w, r)
	})
	handleFile := func(w http.ResponseWriter, r *http.Request) {
		e.handleSpecifiedFile(r.URL.Path, w, r)
	}
	router.HandleFunc("/bootstrap/", handleFile)
	/*
		router.HandleFunc("/fontawesome/", handleFile)
	*/
	router.HandleFunc("/static/", handleFile)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			e.notFound(w, r)
		} else {
			WriteResponse(w, r.URL.Path, e.Pages["home"])
		}
	})
	/*
		router.HandleFunc("/istoriya_prosmotrov", func(w http.ResponseWriter, r *http.Request) {
			WriteResponse(w, r.URL.Path, e.Pages["history"])
		})
		router.HandleFunc("/dostavka", func(w http.ResponseWriter, r *http.Request) {
			WriteResponse(w, r.URL.Path, e.Pages["delivery"])
		})
		router.HandleFunc("/kontakty", func(w http.ResponseWriter, r *http.Request) {
			WriteResponse(w, r.URL.Path, e.Pages["contact"])
		})
		router.HandleFunc("/zakazy/novyy", func(w http.ResponseWriter, r *http.Request) {
			WriteResponse(w, r.URL.Path, e.Pages["neworder"])
		})
	*/
	for _, product := range products {
		magazineURL := "/zhurnaly/" + product.ID
		magazineImageURL := "/product/image/" + product.ID

		func(url string, product *Product) {
			htmlData := product.Thumb
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				WriteResponse(w, r.URL.Path, htmlData)
			})
		}(magazineImageURL, product)

		func(newImgURL string, product *Product) {
			url := fmt.Sprintf("/thumbs/magazine/gallery/%v", product.OldImgName)
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, newImgURL, http.StatusMovedPermanently)
			})
		}(magazineImageURL, product)

		func(url string, product *Product) {
			body := html.Generate("show product", html.Product, product).String()
			htmlData := NewProductPage(body, product.Title).HTMLData()
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				WriteResponse(w, r.URL.Path, htmlData)
			})
		}(magazineURL, product)

		func(newURL string, product *Product) {
			oldURL := fmt.Sprintf("/zhurnaly/%v", product.OldID)
			router.HandleFunc(oldURL, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, newURL, http.StatusMovedPermanently)
			})
		}(magazineURL, product)
	}

	for i, body := range productPagesHtml {
		func(index int, htmlData []byte) {
			url := fmt.Sprintf("/zhurnaly/stranitsa/%v", i+1)
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				WriteResponse(w, r.URL.Path, htmlData)
			})
		}(i, NewProductsPage(body).HTMLData())
	}

	router.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			e.notFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
		defer r.Body.Close()

		foundProducts, err := e.readProducts(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, r.URL.Path, []byte(err.Error()))
			return
		}

		buf := html.Generate("product list", html.ProductList, foundProducts)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		WriteResponse(w, r.URL.Path, buf.Bytes())
	})

	router.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			e.notFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
		defer r.Body.Close()

		foundProducts, err := e.readProducts(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteResponse(w, r.URL.Path, []byte(err.Error()))
			return
		}

		if len(foundProducts) != 0 {
			c := Cart{
				Products: foundProducts,
				Empty:    len(foundProducts) == 0,
			}

			for _, p := range c.Products {
				c.TotalPrice += p.Price
			}

			buf := html.Generate("cart list", html.OrderedProducts, c)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			WriteResponse(w, r.URL.Path, buf.Bytes())
		} else {
			WriteResponse(w, r.URL.Path, []byte{})
		}
	})
	e.router = router
}

func (e *EMarket) readProducts(r io.Reader) ([]*Product, error) {
	reqProducts := make([]string, 0)
	if err := json.NewDecoder(r).Decode(&reqProducts); err != nil {
		return nil, err
	}

	var foundProducts []*Product
	for _, id := range reqProducts {
		if product, found := e.ProductsMap[id]; found {
			foundProducts = append(foundProducts, product)
		}
	}

	return foundProducts, nil
}

func (c *Content) detectType(filename string) string {
	for suffix, contentType := range c.mytype {
		if strings.HasSuffix(filename, suffix) {
			return contentType
		}
	}

	panic(fmt.Sprintf("unknown type %v", filename))
}

func WriteResponse(w http.ResponseWriter, path string, data []byte) {
	if _, err := w.Write(data); err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			fmt.Printf("%v %v\n", path, err)
		}
	}
}

func (e *EMarket) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	WriteResponse(w, r.URL.Path, e.Pages["notfound"])
}

func (e *EMarket) handleSpecifiedFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := readFile(e.rootDir + filename)
	if err == nil {
		w.Header().Set("Content-Type", e.content.detectType(r.URL.Path))
		WriteResponse(w, r.URL.Path, body)
	} else {
		e.notFound(w, r)
	}
}

func (e *EMarket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func readFile(filename string) ([]byte, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("read file: %s %v\n", filename, err)
		return nil, err
	}

	if !stat.Mode().IsRegular() {
		fmt.Printf("not regular file: %s\n", filename)
		return nil, errors.New("not a regulat file")
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("cannot read file %s %v\n", filename, err)
	}
	return body, err
}

func LoadProducts(dataPath string) (products []*Product, err error) {
	data, err := ioutil.ReadFile(dataPath)

	if err != nil {
		return products, err
	}

	var tmpProducts []*Product
	if err := json.Unmarshal(data, &tmpProducts); err != nil {
		return products, err
	}

	for _, product := range tmpProducts {
		if product.Enable {
			products = append(products, product)
		}
	}
	return products, nil
}
