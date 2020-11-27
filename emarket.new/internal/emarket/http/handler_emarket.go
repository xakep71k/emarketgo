package http

import (
	"bytes"
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html/page"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type EMarketHandler struct {
	magazStorage emarket.MagazineStorage
	router       *http.ServeMux
	webRoot      string
	fileCache    map[string][]byte
	notFoundPage []byte
}

func NewEMarketHandler(webRoot string, magazStorage emarket.MagazineStorage) *EMarketHandler {
	h := &EMarketHandler{
		magazStorage: magazStorage,
		router:       http.NewServeMux(),
		webRoot:      webRoot,
		fileCache:    make(map[string][]byte),
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
		magazineImageURL := "/product/image/" + magaz.ID

		func(url string, magaz *emarket.Magazine) {
			data := magaz.Thumb
			e.router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				setCacheControl(w)
				writeResponse(w, r.URL.Path, data)
			})
		}(magazineImageURL, magaz)
	}
}

func concatFiles(rootDir string, files []string) ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, file := range files {
		data, err := ioutil.ReadFile(rootDir + "/" + file)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.Write([]byte("\n"))
	}

	return buf.Bytes(), nil
}
