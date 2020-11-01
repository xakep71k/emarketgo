package mux

import (
	"emarket/internal/magazine"
	"net/http"
)

type Router struct {
	core *http.ServeMux
}

func New(m *magazine.Controller) *Router {
	router := http.NewServeMux()
	return &Router{core: router}
}

func (r *Router) ServeHTTP(http.ResponseWriter, *http.Request) {
}
