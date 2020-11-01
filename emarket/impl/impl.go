package impl

import (
	"emarket/html"
	"emarket/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Cart struct {
	Products   []*model.Product
	Empty      bool
	TotalPrice int
}

const ProductsTitle = "Журналы и выкройки для шитья"

type EMarket struct {
	rootDir     string
	content     *Content
	router      *http.ServeMux
	Pages       map[string][]byte
	ProductsMap map[string]*model.Product
	CSS         []byte
	JS          []byte
	FileCache   map[string][]byte
}

func NewEMarket(rootDir string, db *model.DB) (*EMarket, error) {
	products, err := db.FindAllProducts(model.NewFilter().Enable(true).Sort(true))

	if err != nil {
		return nil, err
	}

	e := &EMarket{
		rootDir:     rootDir,
		content:     NewContent(),
		ProductsMap: make(map[string]*model.Product, 0),
		FileCache:   make(map[string][]byte),
	}

	for _, product := range products {
		e.ProductsMap[product.ID] = product
	}

	productPages := collectProductsPages(products)
	setPageNum(productPages)
	productPagesHtml := buildHtmlProductPages(productPages)
	e.Pages = map[string][]byte{
		"contact":  NewContactPage().htmlData(),
		"delivery": NewDeliveryPage().htmlData(),
		"notfound": NewNotFoundPage().htmlData(),
		"home":     NewHomePage(productPagesHtml[0]).htmlData(),
		"history":  NewHistoryPage().htmlData(),
		"neworder": NewOrderPage().htmlData(),
	}

	allCSS, err := concatFiles(rootDir, html.CSSs)
	if err != nil {
		return nil, err
	}

	allJS, err := concatFiles(rootDir, html.JSs)
	if err != nil {
		return nil, err
	}

	e.FileCache["/static/css/all.css"] = doMinify(allCSS, "text/css")
	e.FileCache["/static/js/all.js"] = doMinify(allJS, "application/javascript")
	e.setupRouter(products, productPagesHtml)
	return e, nil
}

func (e *EMarket) searchFile(file string) string {
	for _, path := range html.SearchPath {
		full := e.rootDir + path + file
		if _, err := os.Stat(full); err == nil {
			return full
		}
	}

	return ""
}

func (e *EMarket) staticHandler(w http.ResponseWriter, r *http.Request) {
	log := func(err error) {
		fmt.Printf("%v %v", r.URL.Path, err)
	}

	requestedFile, err := filepath.Abs(r.URL.Path)
	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	content := e.FileCache[requestedFile]
	if content == nil {
		fullPath := e.searchFile(strings.TrimPrefix(requestedFile, "/static/"))
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
		e.FileCache[requestedFile] = content
	}

	ctype, err := e.content.detectType(requestedFile)
	if err != nil {
		log(err)
		e.notFound(w, r)
		return
	}

	setCacheControl(w)
	w.Header().Set("Content-Type", ctype)
	writeResponse(w, r.URL.Path, content)
}

func (e *EMarket) setupRouter(products []*model.Product, productPagesHtml []string) {
	router := http.NewServeMux()
	router.HandleFunc("/favicon.ico", e.staticHandler)
	router.HandleFunc("/static/", e.staticHandler)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			e.notFound(w, r)
		} else {
			writeResponse(w, r.URL.Path, e.Pages["home"])
		}
	})
	router.HandleFunc("/istoriya_prosmotrov", func(w http.ResponseWriter, r *http.Request) {
		writeResponse(w, r.URL.Path, e.Pages["history"])
	})
	/*
		router.HandleFunc("/dostavka", func(w http.ResponseWriter, r *http.Request) {
			writeResponse(w, r.URL.Path, e.Pages["delivery"])
		})
		router.HandleFunc("/kontakty", func(w http.ResponseWriter, r *http.Request) {
			writeResponse(w, r.URL.Path, e.Pages["contact"])
		})
		router.HandleFunc("/zakazy/novyy", func(w http.ResponseWriter, r *http.Request) {
			writeResponse(w, r.URL.Path, e.Pages["neworder"])
		})
	*/
	for _, product := range products {
		magazineURL := "/zhurnaly/" + product.ID
		magazineImageURL := "/product/image/" + product.ID

		func(url string, product *model.Product) {
			htmlData := product.Thumb
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				setCacheControl(w)
				writeResponse(w, r.URL.Path, htmlData)
			})
		}(magazineImageURL, product)

		func(newImgURL string, product *model.Product) {
			url := fmt.Sprintf("/thumbs/magazine/gallery/%v", product.OldImgName)
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, newImgURL, http.StatusMovedPermanently)
			})
		}(magazineImageURL, product)

		func(url string, product *model.Product) {
			body := html.Generate("show product", html.Product, product).String()
			htmlData := NewProductPage(body, product.Title).htmlData()
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				writeResponse(w, r.URL.Path, htmlData)
			})
		}(magazineURL, product)

		func(newURL string, product *model.Product) {
			oldURL := fmt.Sprintf("/zhurnaly/%v", product.OldID)
			router.HandleFunc(oldURL, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, newURL, http.StatusMovedPermanently)
			})
		}(magazineURL, product)
	}

	for i, body := range productPagesHtml {
		func(index int, htmlData []byte) {
			url := fmt.Sprintf("/zhurnaly/stranitsa/%v", i+1)
			router.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
				writeResponse(w, r.URL.Path, htmlData)
			})
		}(i, NewProductsPage(body).htmlData())
	}

	router.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			e.notFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
		foundProducts, err := e.readProducts(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeResponse(w, r.URL.Path, []byte(err.Error()))
			return
		}

		buf := html.Generate("product list", html.ProductList, foundProducts)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		writeResponse(w, r.URL.Path, buf.Bytes())
	})

	router.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			e.notFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
		foundProducts, err := e.readProducts(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeResponse(w, r.URL.Path, []byte(err.Error()))
			return
		}

		if len(foundProducts) != 0 {
			c := Cart{
				Products: foundProducts,
				Empty:    len(foundProducts) == 0,
			}

			for _, p := range c.Products {
				c.TotalPrice += p.Price
			}

			buf := html.Generate("cart list", html.OrderedProducts, c)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			writeResponse(w, r.URL.Path, buf.Bytes())
		} else {
			writeResponse(w, r.URL.Path, []byte{})
		}
	})
	e.router = router
}

func (e *EMarket) readProducts(r io.Reader) ([]*model.Product, error) {
	reqProducts := make([]string, 0)
	if err := json.NewDecoder(r).Decode(&reqProducts); err != nil {
		return nil, err
	}

	var foundProducts []*model.Product
	for _, id := range reqProducts {
		if product, found := e.ProductsMap[id]; found {
			foundProducts = append(foundProducts, product)
		}
	}

	return foundProducts, nil
}

func (e *EMarket) notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	fmt.Printf("not found %v\n", r.URL.Path)
	writeResponse(w, r.URL.Path, e.Pages["notfound"])
}

func (e *EMarket) handleSpecifiedFile(filename string, w http.ResponseWriter, r *http.Request) {
	body, err := readFile(e.rootDir + filename)
	if err == nil {
		ctype, err := e.content.detectType(r.URL.Path)
		if err == nil {
			w.Header().Set("Content-Type", ctype)
			writeResponse(w, r.URL.Path, body)
		} else {
			fmt.Printf("%v %v\n", r.URL.Path, err)
			e.notFound(w, r)
		}
	} else {
		fmt.Printf("%v %v\n", r.URL.Path, err)
		e.notFound(w, r)
	}
}

func (e *EMarket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.ServeHTTP(w, r)
}
