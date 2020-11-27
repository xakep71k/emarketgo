package http

import (
	"errors"
	"fmt"
	"net/http"
	"syscall"
)

func writeResponse(w http.ResponseWriter, path string, data []byte) {
	if _, err := w.Write(data); err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			fmt.Printf("%v %v\n", path, err)
		}
	}
}
