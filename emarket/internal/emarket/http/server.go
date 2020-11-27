package http

import (
	"net/http"
	"time"
)

type Server = http.Server
type Handler = http.Handler

func NewServer(handler Handler, listenAddr string) *Server {
	return &Server{
		Handler:      handler,
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}
