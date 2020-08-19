package impl

import (
	"bytes"
	"emarket/html"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"syscall"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	mhtml "github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

func LoadProducts(dataPath string) (products []*Product, err error) {
	data, err := ioutil.ReadFile(dataPath)

	if err != nil {
		return products, err
	}

	var tmpProducts []*Product
	if err := json.Unmarshal(data, &tmpProducts); err != nil {
		return products, err
	}

	for _, product := range tmpProducts {
		if product.Enable {
			products = append(products, product)
		}
	}
	return products, nil
}

const pageSize = 30

func buildHtmlProductPages(productPages [][]*Product) []string {
	pages := make([]string, 0)
	maxPages := len(productPages)

	if maxPages == 0 {
		panic("there are no any products")
	}

	for index, products := range productPages {
		plist := html.Generate("product list", html.ProductList, products).String()
		pagination := html.Generate(
			"pagination",
			html.Pagination,
			NewProductPageList(index, maxPages, products, plist),
		).String()
		pages = append(pages, pagination)
	}
	return pages
}

func collectProductsPages(products []*Product) (productPages [][]*Product) {
	iPage := -1
	for i, p := range products {
		if (i % pageSize) == 0 {
			iPage++
			productPages = append(productPages, make([]*Product, 0))
		}

		productPages[iPage] = append(productPages[iPage], p)
	}

	return
}

func setPageNum(productPages [][]*Product) {
	for index, products := range productPages {
		for _, p := range products {
			p.PageNum = index + 1
		}
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

func setCacheControl(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "max-age=31536000")
}

func writeResponse(w http.ResponseWriter, path string, data []byte) {
	if _, err := w.Write(data); err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			fmt.Printf("%v %v\n", path, err)
		}
	}
}

func readFile(filename string) ([]byte, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("read file: %s %v\n", filename, err)
		return nil, err
	}

	if !stat.Mode().IsRegular() {
		fmt.Printf("not regular file: %s\n", filename)
		return nil, errors.New("not a regulat file")
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("cannot read file %s %v\n", filename, err)
	}
	return body, err
}

func doMinify(body []byte, mtype string) []byte {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", mhtml.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	b, err := m.Bytes(mtype, body)
	if err != nil {
		panic("minify " + err.Error())
	}

	return b
}
