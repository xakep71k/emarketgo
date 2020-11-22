package template

import (
	"bytes"
	"text/template"
)

type Page struct {
	Title string
	Body  string
}

type PageBuilder struct {
	template *template.Template
	html     string
	name     string
	args     interface{}
	funcs    map[string]interface{}
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

	page := &Page{
		Body:  buf.String(),
		Title: title,
	}

	return page, nil
}

func NewPageBuilder() *PageBuilder {
	return &PageBuilder{}
}
