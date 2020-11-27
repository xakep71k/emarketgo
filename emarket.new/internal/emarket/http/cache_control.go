package http

import "net/http"

func setCacheControl(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "max-age=31536000")
}
