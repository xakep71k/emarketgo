package http

import (
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html/page"
	"errors"
	"fmt"
	"log"
	"net/http"
	"syscall"
)

type EMarketHandler struct {
	magazStorage emarket.MagazineStorage
	router       *http.ServeMux
}

func NewEMarketHandler(magazStorage emarket.MagazineStorage) *EMarketHandler {
	h := &EMarketHandler{
		magazStorage: magazStorage,
		router:       http.NewServeMux(),
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

func writeResponse(w http.ResponseWriter, path string, data []byte) {
	if _, err := w.Write(data); err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			fmt.Printf("%v %v\n", path, err)
		}
	}
}
