package template

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

type Page struct {
	Title string
	Body  string
}

type PageBuilder struct {
	template    *template.Template
	html        string
	name        string
	args        interface{}
	funcs       map[string]interface{}
	doNotMinify bool
}

func (r *PageBuilder) TemplateFuncs(funcs map[string]interface{}) *PageBuilder {
	r.funcs = funcs
	return r
}

func (r *PageBuilder) Name(name string) *PageBuilder {
	r.name = name
	return r
}

func (r *PageBuilder) Template(html string) *PageBuilder {
	r.html = html
	return r
}

func (r *PageBuilder) Args(args interface{}) *PageBuilder {
	r.args = args
	return r
}

func (r *PageBuilder) Build(title string) (*Page, error) {
	templ, err := template.New(r.name).Funcs(r.funcs).Parse(r.html)

	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	if err := templ.Execute(buf, r.args); err != nil {
		return nil, err
	}

	body := ""

	if r.doNotMinify {
		body = buf.String()
	} else {
		body = string(doMinify(buf.Bytes(), "text/html"))
	}

	page := &Page{
		Body:  body,
		Title: title,
	}

	return page, nil
}

func NewPageBuilder() *PageBuilder {
	return &PageBuilder{}
}

func doMinify(body []byte, mtype string) []byte {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	b, err := m.Bytes(mtype, body)
	if err != nil {
		panic("minify " + err.Error())
	}

	return b
}
