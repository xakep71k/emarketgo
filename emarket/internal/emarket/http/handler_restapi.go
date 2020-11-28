package http

import (
	"emarket/internal/emarket"
	"emarket/internal/emarket/http/internal/html/page"
	"encoding/json"
	"io"
	"net/http"
)

const RequestLimit = 1024 * 2

func (e *EMarketHandler) setupRESTAPI(allMagaz []*emarket.Magazine) {
	magazMap := newMagazMap(allMagaz)

	e.router.HandleFunc("/api/magazines", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			e.notFound(w, r)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
		foundMagazes, err := readMagazineReq(r.Body, magazMap)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeResponse(w, r.URL.Path, []byte(err.Error()))
			return
		}

		content := page.Magazines(foundMagazes)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		writeResponse(w, r.URL.Path, content)
	})

}

func readMagazineReq(r io.Reader, magazMap map[string]*emarket.Magazine) ([]*emarket.Magazine, error) {
	reqMagazes := make([]string, 0)

	if err := json.NewDecoder(r).Decode(&reqMagazes); err != nil {
		return nil, err
	}

	var foundMagazes []*emarket.Magazine

	for _, id := range reqMagazes {
		if magaz, found := magazMap[id]; found {
			foundMagazes = append(foundMagazes, magaz)
		}
	}

	return foundMagazes, nil
}

func newMagazMap(mm []*emarket.Magazine) map[string]*emarket.Magazine {
	magazMap := make(map[string]*emarket.Magazine)

	for _, m := range mm {
		magazMap[m.ID] = m
	}

	return magazMap
}
