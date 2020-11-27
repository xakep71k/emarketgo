package http

import (
	"emarket/internal/pkg/minify"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var CSSs = []string{
	"/static/css/custom_bootstrap.css",
	"/static/css/app.css",
}

var JSs = []string{
	"/static/js/app.js",
}

func (e *EMarketHandler) setupFileHandler() {
	allCSS, err := concatFiles(e.webRoot, CSSs)
	if err != nil {
		log.Fatalln(err)
	}

	allJS, err := concatFiles(e.webRoot, JSs)
	if err != nil {
		log.Fatalln(err)
	}

	e.fileCache["/static/css/all.css"] = minify.DoMinify(allCSS, "text/css")
	e.fileCache["/static/js/all.js"] = minify.DoMinify(allJS, "application/javascript")

	const favicon = "/favicon.ico"
	faviconPath := "/static" + favicon

	e.router.HandleFunc(favicon, func(w http.ResponseWriter, r *http.Request) {
		e.handleSpecifiedFile(w, r, faviconPath)
	})

	e.router.HandleFunc("/static/", e.fileHandler)
}

func (e *EMarketHandler) fullpath(file string) string {
	full := e.webRoot + file

	if _, err := os.Stat(full); err == nil {
		return full
	}

	return ""
}

func (e *EMarketHandler) fileHandler(w http.ResponseWriter, r *http.Request) {
	e.handleSpecifiedFile(w, r, r.URL.Path)
}

func (e *EMarketHandler) handleSpecifiedFile(w http.ResponseWriter, r *http.Request, filename string) {
	log := func(err error) {
		fmt.Printf("%v %v", r.URL.Path, err)
	}

	requestedFile, err := filepath.Abs(filename)

	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	content := e.fileCache[requestedFile]
	ctype, err := detectType(requestedFile)

	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	if content == nil {
		fullPath := e.fullpath(requestedFile)

		if fullPath == "" {
			e.notFound(w, r)
			return
		}

		var err error
		content, err = ioutil.ReadFile(fullPath)

		if err != nil {
			log(err)
			e.notFound(w, r)
			return
		}

		e.fileCache[requestedFile] = minify.DoMinify(content, ctype)
	}

	setCacheControl(w)
	w.Header().Set("Content-Type", ctype)
	writeResponse(w, r.URL.Path, content)
}
