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
}

func NewEMarketHandler(webRoot string, magazStorage emarket.MagazineStorage) *EMarketHandler {
	h := &EMarketHandler{
		magazStorage: magazStorage,
		router:       http.NewServeMux(),
		webRoot:      webRoot,
	}

	h.setupRouter()

	return h
}

func (e *EMarketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}

func (e *EMarketHandler) setupRouter() {
	e.setupMagazPages()
}

func (e *EMarketHandler) setupMagazPages() {
	allMagaz, err := e.magazStorage.Find()

	if err != nil {
		log.Fatalln(err)
	}

	pages := page.MagazineList(allMagaz)

	for i, data := range pages {
		func(index int, htmlData []byte) {
			url := fmt.Sprintf("/zhurnaly/stranitsa/%d", i+1)
			e.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				writeResponse(w, r.URL.Path, htmlData)
			})
		}(i, data)
	}
}
