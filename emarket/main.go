package main

import (
	"emarket/html"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	PageNo = iota
	PageHome
	PageDelivery
	PageContact
)

const (
	URLHome     = "/"
	URLDelivery = "/dostavka"
	URLContact  = "/kontakty"
)

const (
	TitleHome     = "Журналы"
	TitleDelivery = "Доставка"
	TitleContact  = "Контакты"
	TitleNotFound = "Страница не найдена"
)

type EMarketPage struct {
	Title       string
	Body        string
	CurrentPage int
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
	rootDir string
	content *Content
	router  *http.ServeMux
}

func NewEMarket(rootDir string) *EMarket {
	e := &EMarket{
		rootDir: rootDir,
		content: NewContent(),
	}
	router := http.NewServeMux()
	router.HandleFunc("/favicon.ico", e.handleFavicon)
	router.HandleFunc("/bootstrap/", e.handleFile)
	router.HandleFunc("/fontawesome/", e.handleFile)
	router.HandleFunc("/static/", e.handleFile)
	router.HandleFunc(URLHome, e.home)
	router.HandleFunc(URLDelivery, e.Delivery)
	router.HandleFunc(URLContact, e.ContactUs)
	e.router = router
	return e
}

func (c *Content) DetectType(filename string) string {
	for suffix, contentType := range c.mytype {
		if strings.HasSuffix(filename, suffix) {
			return contentType
		}
	}

	panic(fmt.Sprintf("unknown type %v", filename))
}

func (e *EMarket) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	p := EMarketPage{Title: TitleNotFound, Body: html.NotFound}
	html.AppPage.Execute(w, p)
}

func (e *EMarket) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		e.notFound(w, r)
	} else {
		p := EMarketPage{Title: TitleHome, CurrentPage: PageHome}
		html.AppPage.Execute(w, p)
	}
}

func (e *EMarket) handleCustomFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := readFile(e.rootDir + r.URL.Path)
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

func (e *EMarket) Delivery(w http.ResponseWriter, r *http.Request) {
	p := EMarketPage{Title: TitleDelivery, CurrentPage: PageDelivery}
	html.AppPage.Execute(w, p)
}

func (e *EMarket) ContactUs(w http.ResponseWriter, r *http.Request) {
	p := EMarketPage{Title: TitleContact, CurrentPage: PageContact}
	html.AppPage.Execute(w, p)
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

	emarket := NewEMarket(webRoot)
	srv := &http.Server{
		Handler:      emarket,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
