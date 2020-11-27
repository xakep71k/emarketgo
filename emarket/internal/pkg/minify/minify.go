package minify

import (
	"log"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

var minif = minify.New()
var types map[string]minify.MinifierFunc

func init() {
	types = make(map[string]minify.MinifierFunc)
	types["text/css"] = css.Minify
	types["text/html"] = html.Minify
	types["application/javascript"] = js.Minify

	for ctype, minifier := range types {
		minif.AddFunc(ctype, minifier)
	}
}

func DoMinify(body []byte, mtype string) []byte {
	if _, found := types[mtype]; !found {
		return body
	}

	b, err := minif.Bytes(mtype, body)

	if err != nil {
		log.Printf("minify: %s\n", err.Error())
		return body
	}

	return b
}
