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

type WEBPage struct {
	Title       string
	Body        string
	CurrentPage int
}

var WEBRoot string
var ContentTypeHash map[string]string

func init() {
	ContentTypeHash = map[string]string{
		".js":     "application/javascript",
		".css":    "text/css",
		".ico":    "image/x-icon",
		"js.map":  "application/octet-stream",
		"css.map": "application/octet-stream",
		".woff":   "font/woff",
		".woff2":  "font/woff2",
		".ttf":    "font/ttf",
	}
}

func ContentType(filename string) string {
	for suffix, contentType := range ContentTypeHash {
		if strings.HasSuffix(filename, suffix) {
			return contentType
		}
	}

	panic(fmt.Sprintf("unknown type %v", filename))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	p := WEBPage{Title: "Страница не найдена", Body: html.NotFound}
	html.AppPage.Execute(w, p)
}

func Magazine(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "inside magazine")
}

func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFound(w, r)
	} else {
		p := WEBPage{Title: "test title", CurrentPage: PageHome}
		html.AppPage.Execute(w, p)
	}
}

func HandleFile(w http.ResponseWriter, r *http.Request) {
	body, err := readFile(r.URL.Path)
	if err == nil {
		w.Header().Set("Content-Type", ContentType(r.URL.Path))
		w.Write(body)
	} else {
		NotFound(w, r)
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	body, err := readFile("/static/favicon.ico")
	if err == nil {
		w.Header().Set("Content-Type", ContentType(r.URL.Path))
		w.Write(body)
	}
}

func Delivery(w http.ResponseWriter, r *http.Request) {
	p := WEBPage{Title: "Доставка", CurrentPage: PageDelivery}
	html.AppPage.Execute(w, p)
}

func ContactUs(w http.ResponseWriter, r *http.Request) {
	p := WEBPage{Title: "Контакты", CurrentPage: PageContact}
	html.AppPage.Execute(w, p)
}

func readFile(filename string) ([]byte, error) {
	filename = WEBRoot + filename
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
	var err error
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s --web-root <path>\n", os.Args[0])
		os.Exit(1)
	}
	webRoot := flag.String("web-root", "", "<path>")
	flag.Parse()

	if webRoot == nil || *webRoot == "" {
		fmt.Println("web root not specified")
		os.Exit(1)
	}

	WEBRoot, err = filepath.Abs(*webRoot)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	router := http.NewServeMux()
	router.HandleFunc("/favicon.ico", Favicon)
	router.HandleFunc("/bootstrap/", HandleFile)
	router.HandleFunc("/fontawesome/", HandleFile)
	router.HandleFunc("/static/", HandleFile)
	router.HandleFunc("/", Root)
	router.HandleFunc("/dostavka", Delivery)
	router.HandleFunc("/kontakty", ContactUs)
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
