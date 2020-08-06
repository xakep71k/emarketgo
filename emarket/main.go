package main

import (
	"bytes"
	"emarket/db"
	"emarket/html"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	PageNo = iota
	PageHome
	PageDelivery
	PageContact
)

const pageSize = 30

type Product struct {
	ID          string `bson:"_id,omitempty"`
	Title       string `bson:"title"`
	Price       int    `bson:"price"`
	Thumb       []byte `bson:"thumb"`
	Enable      bool   `bson:"enable"`
	Description string `bson:"description"`
	Quantity    int    `bson:"quantity"`
	OldID       int    `bson:"oldid"`
	OldImgName  string `bson:"oldimgfile"`
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
	Index       int
	Products    []*Product
	PageNumbers []int
	MaxPages    int
	ListHTML    string
}

func NewProductPageList(index int, maxPages int, products []*Product, listHtml string) *ProductPageList {
	p := &ProductPageList{
		Index:    index,
		First:    index == 0,
		Last:     maxPages == 0 || maxPages == index+1,
		Products: products,
		MaxPages: maxPages,
		ListHTML: listHtml,
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
		Title:       "Журналы",
		CurrentPage: PageHome,
		Body:        body,
	}

	return p
}

func NewDeliveryPage() *Page {
	return &Page{
		ID:          "service-delivery_terms",
		Title:       "Доставка",
		CurrentPage: PageDelivery,
		Body:        html.Delivery,
	}
}

func NewContactPage() *Page {
	return &Page{
		Title:       "Контакты",
		CurrentPage: PageContact,
		Body:        html.Contact,
	}
}

func NewNotFoundPage() *Page {
	return &Page{
		Title: "Страница не найдена",
		Body:  html.NotFound,
	}
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
		"js.map":  "application/octet-stream",
		"css.map": "application/octet-stream",
		".woff":   "font/woff",
		".woff2":  "font/woff2",
		".ttf":    "font/ttf",
	}
	return c
}

type EMarket struct {
	rootDir              string
	content              *Content
	router               *http.ServeMux
	Pages                map[string]*Page
	productsIDHash       map[string]*Product
	productsFileNameHash map[string]*Product
	productsOldIDHash    map[string]*Product
	productPagesHtml     []string
	productPages         [][]*Product
}

func buildHtmlProductPages(productPages [][]*Product) []string {
	pages := make([]string, 0)
	funcMap := template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	}
	generate := func(name, templateStr string, input interface{}) string {
		templ, err := template.New(name).Funcs(funcMap).Parse(templateStr)
		if err != nil {
			panic(fmt.Sprintf("%v %v", name, err))
		}

		var buf bytes.Buffer
		if err := templ.Execute(&buf, input); err != nil {
			panic(fmt.Sprintf("%v %v", name, err))
		}

		return buf.String()
	}

	maxPages := len(productPages)
	if maxPages == 0 {
		panic("there are no any products")
	}
	for index, products := range productPages {
		plist := generate("product list", html.ProductList, products)
		pagination := generate("pagination", html.Pagination, NewProductPageList(index, maxPages, products, plist))
		pages = append(pages, pagination)
	}
	return pages
}

func NewEMarket(rootDir string, products []*Product) (*EMarket, error) {
	sort.SliceStable(products, func(i, j int) bool {
		return products[i].Title < products[j].Title
	})

	e := &EMarket{
		rootDir:              rootDir,
		content:              NewContent(),
		productsIDHash:       make(map[string]*Product),
		productsOldIDHash:    make(map[string]*Product),
		productsFileNameHash: make(map[string]*Product),
	}

	var productPages [][]*Product
	iPage := -1
	for i, p := range products {
		if (i % pageSize) == 0 {
			iPage++
			productPages = append(productPages, make([]*Product, 0))
		}

		productPages[iPage] = append(productPages[iPage], p)
		e.productsIDHash[p.ID] = p
		e.productsOldIDHash[fmt.Sprintf("%v", p.OldID)] = p
		e.productsFileNameHash[p.OldImgName] = p
	}

	e.productPages = productPages
	e.productPagesHtml = buildHtmlProductPages(productPages)

	e.Pages = map[string]*Page{
		"contact":  NewContactPage(),
		"delivery": NewDeliveryPage(),
		"notfound": NewNotFoundPage(),
		"home":     NewHomePage(e.productPagesHtml[0]),
	}

	e.setupRouter()
	return e, nil
}

func (e *EMarket) setupRouter() {
	router := http.NewServeMux()
	router.HandleFunc("/favicon.ico", e.handleFavicon)
	router.HandleFunc("/bootstrap/", e.handleFile)
	router.HandleFunc("/fontawesome/", e.handleFile)
	router.HandleFunc("/static/", e.handleFile)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			e.notFound(w, r)
		} else {
			e.home(w, r)
		}
	})
	router.HandleFunc("/dostavka", e.delivery)
	router.HandleFunc("/kontakty", e.contact)
	for id := range e.productsIDHash {
		func(key string) {
			router.HandleFunc("/product/image/"+key, func(w http.ResponseWriter, r *http.Request) {
				w.Write(e.productsIDHash[key].Thumb)
			})
		}(id)
	}

	for file := range e.productsFileNameHash {
		func(key string) {
			router.HandleFunc("/thumbs/magazine/gallery/"+key, func(w http.ResponseWriter, r *http.Request) {
				w.Write(e.productsFileNameHash[key].Thumb)
			})
		}(file)
	}

	for i, page := range e.productPagesHtml {
		func(index int, page string) {
			url := fmt.Sprintf("/zhurnaly/stranitsa/%v", i+1)
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				p := &Page{
					ID:          "magazines-index",
					Title:       "Журналы",
					CurrentPage: PageHome,
					Body:        page,
				}
				html.AppPage.Execute(w, p)
			})
		}(i, page)
	}

	e.router = router
}

func (c *Content) DetectType(filename string) string {
	for suffix, contentType := range c.mytype {
		if strings.HasSuffix(filename, suffix) {
			return contentType
		}
	}

	panic(fmt.Sprintf("unknown type %v", filename))
}

func loadProducts(mongoDataPath string) ([]*Product, error) {
	/*
		mongo := db.NewDockerMongo(mongoDataPath)

		if err := mongo.Start(); err != nil {
			return nil, err
		}

		defer mongo.Stop()
	*/

	client, err := db.NewMongoClient()

	if err != nil {
		return nil, err
	}

	defer client.Disconnect(nil)

	collection := client.Database("emarket").Collection("products")
	ctx := db.DefaultContext()
	cursor, err := collection.Find(ctx, bson.M{"enable": true})

	if err != nil {
		return nil, err
	}

	products := make([]*Product, 0)
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (e *EMarket) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	html.AppPage.Execute(w, e.Pages["notfound"])
}

func (e *EMarket) home(w http.ResponseWriter, r *http.Request) {
	html.AppPage.Execute(w, e.Pages["home"])
}

func (e *EMarket) handleSpecifiedFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := readFile(e.rootDir + filename)
	if err == nil {
		w.Header().Set("Content-Type", e.content.DetectType(r.URL.Path))
		w.Write(body)
	} else {
		e.notFound(w, r)
	}
}

func (e *EMarket) handleFile(w http.ResponseWriter, r *http.Request) {
	e.handleSpecifiedFile(r.URL.Path, w, r)
}

func (e *EMarket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *EMarket) handleFavicon(w http.ResponseWriter, r *http.Request) {
	e.handleSpecifiedFile("/static/favicon.ico", w, r)
}

func (e *EMarket) delivery(w http.ResponseWriter, r *http.Request) {
	html.AppPage.Execute(w, e.Pages["delivery"])
}

func (e *EMarket) contact(w http.ResponseWriter, r *http.Request) {
	html.AppPage.Execute(w, e.Pages["contact"])
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

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s --web-root <path> --listen <ip:port> --mongo-data <path>\n", os.Args[0])
		os.Exit(1)
	}
	webRootOpt := flag.String("web-root", "", "<path>")
	listenOpt := flag.String("listen", "", "<ip:port>")
	mongoDataOpt := flag.String("mongo-data", "", "<path>")
	flag.Parse()

	if webRootOpt == nil || *webRootOpt == "" {
		fmt.Println("web root not specified")
		os.Exit(1)
	}

	if listenOpt == nil || *listenOpt == "" {
		fmt.Println("listen ip:port not specified")
		os.Exit(1)
	}

	if mongoDataOpt == nil || *mongoDataOpt == "" {
		fmt.Println("listen ip:port not specified")
		os.Exit(1)
	}

	webRoot, err := filepath.Abs(*webRootOpt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	products, err := loadProducts(*mongoDataOpt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	emarket, err := NewEMarket(webRoot, products)

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      emarket,
		Addr:         *listenOpt,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Println("started")
	log.Fatal(srv.ListenAndServe())
}
