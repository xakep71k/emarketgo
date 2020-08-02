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

func NewHomePage(body string) *Page {
	p := &Page{
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
	productPages         []string
}

func buildProductPages(productPages [][]*Product) []string {
	pages := make([]string, 0)
	for _, page := range productPages {
		templ, err := template.New("product page").Parse(html.ProductList)
		if err != nil {
			panic(err)
		}

		var tpl bytes.Buffer
		if err := templ.Execute(&tpl, page); err != nil {
			panic(err)
		}
		pages = append(pages, tpl.String())
	}
	return pages
}

func NewEMarket(rootDir string) (*EMarket, error) {
	products, err := loadProducts()

	if err != nil {
		return nil, err
	}

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

	e.productPages = buildProductPages(productPages)

	e.Pages = map[string]*Page{
		"contact":  NewContactPage(),
		"delivery": NewDeliveryPage(),
		"notfound": NewNotFoundPage(),
		"home":     NewHomePage(e.productPages[0]),
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

func loadProducts() ([]*Product, error) {
	client, err := db.NewMongoClient()

	if err != nil {
		return nil, err
	}

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

func (e *EMarket) handleCustomFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := readFile(e.rootDir + filename)
	if err == nil {
		w.Header().Set("Content-Type", e.content.DetectType(r.URL.Path))
		w.Write(body)
	} else {
		e.notFound(w, r)
	}
}

func (e *EMarket) handleFile(w http.ResponseWriter, r *http.Request) {
	e.handleCustomFile(r.URL.Path, w, r)
}

func (e *EMarket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *EMarket) handleFavicon(w http.ResponseWriter, r *http.Request) {
	e.handleCustomFile("/static/favicon.ico", w, r)
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
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s --web-root <path>\n", os.Args[0])
		os.Exit(1)
	}
	webRootOpt := flag.String("web-root", "", "<path>")
	flag.Parse()

	if webRootOpt == nil || *webRootOpt == "" {
		fmt.Println("web root not specified")
		os.Exit(1)
	}

	webRoot, err := filepath.Abs(*webRootOpt)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	emarket, err := NewEMarket(webRoot)

	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler:      emarket,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
