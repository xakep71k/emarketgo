package http

import (
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html/page"
	"fmt"
	"log"
	"net/http"
)

type EMarketHandler struct {
	magazStorage emarket.MagazineStorage
	router       *http.ServeMux
	webRoot      string
	notFoundPage []byte
}

func NewEMarketHandler(webRoot string, magazStorage emarket.MagazineStorage) *EMarketHandler {
	h := &EMarketHandler{
		magazStorage: magazStorage,
		router:       http.NewServeMux(),
		webRoot:      webRoot,
		notFoundPage: page.NotFound(),
	}

	h.setupRouter()

	return h
}

func (e *EMarketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *EMarketHandler) setupRouter() {
	allMagaz, err := e.magazStorage.Find()

	if err != nil {
		log.Fatalln(err)
	}

	pages := page.MagazineList(allMagaz)
	e.setupRootPage(pages[0])
	e.setupMagazPages(allMagaz, pages)
	e.setupFileHandler()
	e.setupContactPage()
	e.setupHistoryPage()
	e.setupRESTAPI(allMagaz)
}

func (e *EMarketHandler) setupContactPage() {
	page := page.Contact()
	e.router.HandleFunc("/kontakty", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, r.URL.Path, page)
	})
}

func (e *EMarketHandler) setupHistoryPage() {
	page := page.History()
	e.router.HandleFunc("/istoriya_prosmotrov", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, r.URL.Path, page)
	})
}

func (e *EMarketHandler) setupRootPage(page []byte) {
	e.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			e.notFound(w, r)
		} else {
			writeResponse(w, r.URL.Path, page)
		}
	})
}

func (e *EMarketHandler) setupMagazPages(allMagaz []*emarket.Magazine, pages [][]byte) {
	for i, data := range pages {
		func(index int, htmlData []byte) {
			url := fmt.Sprintf("/zhurnaly/stranitsa/%d", i+1)
			e.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				writeResponse(w, r.URL.Path, htmlData)
			})
		}(i, data)
	}

	for _, magaz := range allMagaz {
		magazineURL := "/zhurnaly/" + magaz.ID
		magazineImageURL := "/product/image/" + magaz.ID

		func(url string, magaz *emarket.Magazine) {
			data := magaz.Thumb
			e.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				setCacheControl(w)
				writeResponse(w, r.URL.Path, data)
			})
		}(magazineImageURL, magaz)

		func(url string, magaz *emarket.Magazine) {
			page := page.Magazine(magaz)
			e.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				writeResponse(w, r.URL.Path, page)
			})
		}(magazineURL, magaz)
	}
}
