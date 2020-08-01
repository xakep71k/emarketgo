package html

import (
	"io"
	"text/template"
)

type htmlTemplate struct {
	template *template.Template
}

func (r *htmlTemplate) Execute(wr io.Writer, data interface{}) {
	if err := r.template.Execute(wr, data); err != nil {
		panic(err)
	}
}
