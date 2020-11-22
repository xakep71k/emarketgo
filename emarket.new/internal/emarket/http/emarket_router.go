package http

import "net/http"

func NewEMarketHandler() *EMarketHandler {
	return &EMarketHandler{}
}

type EMarketHandler struct {
}

func (r *EMarketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.router.ServeHTTP(w, r)
}
