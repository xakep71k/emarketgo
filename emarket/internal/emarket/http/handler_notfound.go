package http

import (
	"fmt"
	"net/http"
)

func (e *EMarketHandler) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	writeResponse(w, r.URL.Path, e.notFoundPage)
}
